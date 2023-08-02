package core

import (
	"database/sql"
	"dataMiner/models"
	"dataMiner/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
	"sort"
	"strings"
)

/*
  Extract data from each table in the database
  @Param  db (database handle)
  @Param  tableList (all the tables in the database)
  @Param  num (the number of rows returned from database)
  @Param  outputID (the output file name)
  @Param  typeD (the type of database)
*/
func Sampledata(db *sql.DB,tableList []string,num int,outputID utils.InfoStruct,typeD string) {
	fmt.Println("Task is in processing......")
	var csv []models.SampleStruct //save data for csv output

	for _, tbl := range tableList {

		val := strings.Split(tbl, ".")
		columnRows:=QueryWrapped(db ,typeD,"column",val[0],val[1],0)
		defer columnRows.Close()

		// Get data from each table
		data:=QueryWrapped(db ,typeD,"data",val[0],val[1],num)
		defer data.Close()
		cols, err := data.Columns()
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Make a slice for the values
		values := make([]sql.RawBytes, len(cols))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		//put into struct
		var ctmp models.SampleStruct
		ctmp.DatabaseName = val[0]
		ctmp.TableName = val[1]

		//Loop through the rows and append the column names to the columnNames
		for columnRows.Next() {
			var columnName string
			if err := columnRows.Scan(&columnName);
				err != nil {log.Fatal(err) }
			ctmp.ColumnName = append(ctmp.ColumnName, columnName)
		}

		//Loop through the rows and append the data to Rows
		for data.Next() {
			// read each row on the table
			// each column value will be stored in the slice
			err = data.Scan(scanArgs...)
			if err != nil {
				log.Fatal("Error scanning rows from table", err)
			}

			var value string
			var line []string

			for _, col := range values {
				// Here  check if the value is nil (NULL value)
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				line = append(line, value)
			}

			//put it into csv
			ctmp.Rows = append(ctmp.Rows,line)
		}
		csv = append(csv, ctmp)
	}
	utils.Savetocsv(csv, outputID,num)
	utils.Savetohtml(csv, outputID)
}

/*
  Extract data from each collection in the database for mongodb
  @Param  client (database handle)
  @Param  collectionList (all the collections in the database)
  @Param  num (the number of rows returned from database)
  @Param  outputID (the output file name)
*/
func SampledataMongo(client *mongo.Client,collectionList []string,num int,outputID utils.InfoStruct) {
	fmt.Println("Task is in processing......")
	var csv []models.SampleStruct //put data into SampleStrut for later output

	for _, clt := range collectionList {

		val := strings.Split(clt, ".")
		var results []bson.M
		GetDocuments(client,val[0],val[1],num,&results)

		//put into struct for output
		var ctmp models.SampleStruct
		ctmp.DatabaseName = val[0]
		ctmp.TableName = val[1]

        //get the documents from each collection
		for counter, result := range results {
			iter := reflect.ValueOf(result).MapRange()
			var docs []models.Document
			for iter.Next() {
				key := iter.Key().String()
				value := iter.Value().Interface()
				var re string
				if value!=""{
					DealWithMongoData(key+".",value,&re)
					re=strings.Trim(re,"\n:")
				}else{
					re=key+": NULL"
				}
				doc:=models.Document{Key:key,Value: re}
				docs=append(docs,doc)
			}
			//order the keys
			sort.Sort(models.Documents(docs))
			var lineSorted []string
			for _,i:=range docs{
				lineSorted=append(lineSorted,i.Value)
			}
			if counter==0{
				for _,i:=range docs{
					ctmp.ColumnName=append(ctmp.ColumnName,i.Key)
				}
			}
			//put the key and data into struct
			ctmp.Rows = append(ctmp.Rows,lineSorted)
		}
		csv = append(csv, ctmp)
	}
	utils.Savetocsv(csv, outputID,num)
	utils.Savetohtml(csv, outputID)
}


/*
  Get the key and data from one document and turn the value into specified format
  @Param  prefix (the key of document)
  @Param  doc (the value of document)
  @Param  re (the document value in specified format)
*/
func DealWithMongoData(prefix string, doc interface{},re *string) {
	// Get the type of the value
	value := reflect.ValueOf(doc)
	switch value.Kind() {
	case reflect.Map:
		// Iterate over each key-value pair in the map
		iter := reflect.ValueOf(doc).MapRange()
		for iter.Next() {
			key := iter.Key().String()
			value := iter.Value().Interface()

			// Recursively get the keys in the value
			DealWithMongoData(prefix+key+".", value,re)
		}
	case reflect.Slice:
		// Iterate over each element in the slice
		for i := 0; i < value.Len(); i++ {
			// Recursively get the keys in the element
			DealWithMongoData(prefix+fmt.Sprintf("[%d]", i)+".", value.Index(i).Interface(),re)
		}
	default:
		prefix=strings.Trim(prefix,".")
		*re=*re+"\n"+prefix+" : "+fmt.Sprint(value)
	}
}