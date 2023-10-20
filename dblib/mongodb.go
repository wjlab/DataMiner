package dblib

import (
	"context"
	"log"
	"strings"
	"time"
	"dataMiner/models"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

/*
  Mongo database initialization, and return database handle (support MongoDB 3.6 and higher)
  @Param  info (the information user inputs)
  @Return *mongo.Client (database handle)
*/
func MongodbInit(info models.InitData) *mongo.Client {
	// Set timeout to 10 seconds
	timeout := time.Second * 10

	// Create a context with the timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Set up a MongoDB client
	clientOptions := options.Client()
	if info.DatabaseUser!=""{
		credential := options.Credential{
			Username: info.DatabaseUser,
			Password: info.DatabasePassword,
			AuthSource: info.AuthSource,
		}
		clientOptions.SetAuth(credential)
	}
	clientOptions.ApplyURI("mongodb://"+info.DatabaseAddress)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	color.Infoln(info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress," connection inited successfully.")
	return client
}

/*
  Mongo: count all collections and add them into list
  @Param  client (database handle)
  @Return []string (all the tables in the database)
*/
func CountAllCollections(client *mongo.Client)([]string){
	var collectionList []string
	// Get the list of database names
	dbNames, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		if strings.Contains(err.Error(),"unable to authenticate using mechanism"){
			log.Fatal("This Mongodb needs to authenticate database name, please provide database name after database address, like: 127.0.0.1:27017?databaseName")
		}else{
			log.Fatal(err)
		}
	}

	for _, dbName := range dbNames {
		if dbName=="config"||dbName=="admin"||dbName=="local"{
			continue
		}
		// Get the list of collection names in each database
		collNames, err := client.Database(dbName).ListCollectionNames(context.Background(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		for _, collName := range collNames {
			collectionList=append(collectionList,dbName+"."+collName)
		}
	}
	return collectionList
}

/*
  Mongo: get all documents from the specified collection and put them into results
  @Param  client (database handle)
  @Param  database (the specified database name)
  @Param  collection (the specified collection name)
  @Param  num (the number of data returned from database)
  @Param  results (save the data returned from database)
*/
func GetDocuments(client *mongo.Client,database,collectionName string,num int ,results *[]bson.M){
	// Select a collection
	collection := client.Database(database).Collection(collectionName)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"_id", 1}})
	//don't show _id object
	findOptions.SetProjection(bson.M{"_id": 0})
	findOptions.SetLimit(int64(num))
	// Get a cursor over all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.D{},findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	if err = cursor.All(context.Background(), results); err != nil {
		log.Fatal(err)
	}
}
