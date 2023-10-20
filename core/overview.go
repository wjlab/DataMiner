package core

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"database/sql"
	"dataMiner/models"
	"dataMiner/utils"
	"dataMiner/dblib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
  Count the amount of data in the database
  @Param  sql.DB (database handle)
  @Param  outputID (the output file name)
  @Param  typeD (the type of database)
*/
func Overview(db *sql.DB,outputID utils.InfoStruct,typeD string) {
	fmt.Println("Task is in processing......")
	var csv []models.OverviewData

	// get all the databases
	dbs:= dblib.QueryWrapped(db,typeD,"database","","",0)
	for dbs.Next() {
		var dbsname string
		err := dbs.Scan(&dbsname)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if typeD == "mysql" {
			if dbsname == "information_schema" || dbsname == "mysql" || dbsname == "performance_schema" || dbsname == "sys" {
				continue
			}
		}else if typeD=="mssql"{
			if dbsname=="tempdb"||dbsname=="master"||dbsname=="model"||dbsname=="msdb"||dbsname=="ReportServer"||dbsname=="ReportServerTempDB"{
				continue
			}
		}else if typeD=="oracle"{
			if dbsname=="XDB"||dbsname=="CTXSYS"||dbsname=="DBSNMP"||dbsname=="EXFSYS"||dbsname=="MDSYS"||dbsname=="ORDSYS"||dbsname=="OLAPSYS"||dbsname=="SYS"||dbsname=="SYSMAN"||dbsname=="FLOWS_FILES"||dbsname=="APEX_030200"||dbsname=="APPQOSSYS"||dbsname=="PM"||dbsname=="ORDDATA"||dbsname=="IX"||dbsname=="WMSYS"||dbsname=="OWBSYS"||dbsname=="OUTLN"||dbsname=="SYSTEM"||dbsname=="OE"||dbsname=="SH"||dbsname=="HR"{
				continue
			}
		}
		// get the tables
		tables:= dblib.QueryWrapped(db,typeD,"table",dbsname,"",0)
		var tblname string
		for tables.Next() {
			err = tables.Scan(&tblname)
			if err != nil {
				log.Fatalf(err.Error())
			}
			var rowCount int
			dblib.QueryCount(db,typeD,dbsname,tblname,&rowCount)
			ctmp:=models.OverviewData{DatabaseName: dbsname,TableName: tblname,RowCount: strconv.Itoa(rowCount)}
			csv=append(csv,ctmp)
		}
	}
	utils.SavetocsvO(csv,outputID)
	utils.SavetohtmlO(csv,outputID)
}

/*
  Count the amount of data in mongo database
  @Param  client (database handle)
  @Param  outputID (the output file name)
*/
func OverviewMongo(client *mongo.Client,outputID utils.InfoStruct){
	fmt.Println("Task is in processing......")
	var csv []models.OverviewData

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
			// Get the number of documents in each collection
			count, err := client.Database(dbName).Collection(collName).CountDocuments(context.Background(), bson.M{})
			if err != nil {
				log.Fatal(err)
			}

			ctmp:=models.OverviewData{DatabaseName: dbName,TableName:collName,RowCount: strconv.Itoa(int(count))}
			csv=append(csv,ctmp)

		}
	}
	utils.SavetocsvO(csv,outputID)
	utils.SavetohtmlO(csv,outputID)
}

/*
  Count the amount of data in Postgre database
  @Param  connectionString (postgre connection string)
  @Param  tableList (all the tables in the database)
  @Param  outputID (the output file name)
*/
func OverviewPostgre(connectionString string,tableList []string,outputID utils.InfoStruct) {
	fmt.Println("Task is in processing......")
	var csv []models.OverviewData //save data for csv output

	for _, tbl := range tableList {
		parts := strings.Split(tbl, ".")
		database := parts[0]
		schema := parts[1]
		table := strings.Join(parts[2:], ".")

		query := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\".\"%s\"",schema, table)
		var rowCount int
		db:=dblib.PostgreDBConnect(connectionString,database)
		err := db.QueryRow(query).Scan(&rowCount)
		if err != nil {
			log.Fatal(err)
		}
		db.Close()
		ctmp:=models.OverviewData{DatabaseName: database,TableName: schema+"."+table,RowCount: strconv.Itoa(rowCount)}
		csv=append(csv,ctmp)
	}
	utils.SavetocsvO(csv,outputID)
	utils.SavetohtmlO(csv,outputID)
}