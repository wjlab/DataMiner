package dblib

import (
	"fmt"
	"log"
	"strings"
	"database/sql"
	"dataMiner/models"
	"github.com/gookit/color"
	_ "github.com/lib/pq"
)

/*
  Postgre database initialization, and return database handle and connection string
  @Param  info (the information user inputs)
  @Return sql.DB (database handle)
  @Return string (connection string)
*/
func PostgreDBinit(info models.InitData) (*sql.DB,string){
	var db *sql.DB
	if info.DatabaseAddress==""{
		log.Fatalf("Enter the database address,please!")
	}
	hostInfo:=strings.Split(info.DatabaseAddress,":")
	if len(hostInfo)!=2{
		log.Fatalf("Please enter the database address in the standard format,eg: 127.0.0.1:5432")
	}
	host:=hostInfo[0]
	port:=hostInfo[1]
	connectionString:= fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port,info.DatabaseUser,info.DatabasePassword)

	db=PostgreDBConnect(connectionString,"")
	err := db.Ping()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	color.Infoln(info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress," connection inited successfully.")
	return db,connectionString
}

/*
  Connect to postgre database and return database handle and connection string
  @Param  connectionString (postgre connection string)
  @Param  dbName (the database name to connect)
  @Return sql.DB (database handle)
*/
func PostgreDBConnect(connectionString,dbName string)(*sql.DB){
	// Create a new connection to the current database
	db, err := sql.Open("postgres", fmt.Sprintf("%s dbname=%s", connectionString, dbName))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

/*
  Count all tables and add them into list
  @Param  db (database handle)
  @Param  connectionString (postgre connection string)
  @Return []string (all the tables in the database)
*/
func CountAllTablesPs(db *sql.DB,connectionString string) ([]string){
	var tableList []string
	// Query to retrieve all databases
	databasesQuery := "SELECT datname FROM pg_database WHERE datistemplate = false"

	// Execute the query to get all databases
	databasesRows, err := db.Query(databasesQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer databasesRows.Close()

	// Iterate over the databases
	for databasesRows.Next() {
		var dbName string
		err = databasesRows.Scan(&dbName)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new connection to the current database
		dbSub, err := sql.Open("postgres", fmt.Sprintf("%s dbname=%s", connectionString, dbName))
		if err != nil {
			log.Fatal(err)
		}
		defer dbSub.Close()

		// Query to retrieve all table names in the current database
		query := `
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_type = 'BASE TABLE'
			AND table_schema NOT LIKE 'pg_%'
			AND table_schema != 'information_schema'
		`

		// Execute the query for the current database
		rows, err := dbSub.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Iterate over the result rows
		for rows.Next() {
			var schemaName, tableName string
			err = rows.Scan(&schemaName, &tableName)
			if err != nil {
				log.Fatal(err)
			}
			tableList=append(tableList,dbName+"."+schemaName+"."+tableName)
		}

		// Check for any errors during iteration
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Check for any errors during iteration
	if err = databasesRows.Err(); err != nil {
		log.Fatal(err)
	}
	return tableList
}

/*
  Count all columns in database and return the corresponding information
  @Param  connectionString (postgre connection string)
  @Param  database (database name)
  @Param  schema (database schema)
  @Param  table (database table)
  @Return *sql.Rows (the information from database)
*/
func QueryColumnsPs(connectionString,database,schema,table string) *sql.Rows{
	db:=PostgreDBConnect(connectionString,database)
	defer db.Close()

	// Query to retrieve all column names in the current table
	columnsQuery := fmt.Sprintf(`
				SELECT column_name
				FROM information_schema.columns
				WHERE table_schema = '%s'
				AND table_name = '%s'
			`, schema, table)

	// Execute the query for the current table
	columnsRows, err := db.Query(columnsQuery)
	if err != nil {
		log.Fatal(err)
	}
	return columnsRows
}

/*
  Get data from the specified database table and return the corresponding information
  @Param  connectionString (postgre connection string)
  @Param  database (database name)
  @Param  schema (database schema)
  @Param  table (database table)
  @Param  num (the number of rows returned from database)
  @Return *sql.Rows (the information from database)
*/
func QueryDataPs(connectionString,database,schema,table string,num int) *sql.Rows{
	db:=PostgreDBConnect(connectionString,database)
	defer db.Close()

	// Query to retrieve the specified rows data in the current table
	dataQuery := fmt.Sprintf("SELECT * FROM \"%s\".\"%s\" LIMIT %d", schema, table,num)

	// Execute the query for the current table
	dataRows, err := db.Query(dataQuery)
	if err != nil {
		log.Fatal(err)
	}

	return dataRows
}