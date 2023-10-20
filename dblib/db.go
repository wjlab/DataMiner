package  dblib

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"database/sql"
	"dataMiner/models"
	"dataMiner/utils"
	"github.com/gookit/color"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/sijms/go-ora/v2"
)

/*
  Mysql database initialization, and return database handle
  @Param  info (the information user inputs)
  @Return sql.DB (database handle)
*/
func MysqlDBinit(info *models.InitData) (*sql.DB){
	var db *sql.DB
	if info.DatabaseAddress==""{
		log.Fatalf("Enter the database address,please!")
	}
	db, err := sql.Open("mysql", info.DatabaseUser+":"+info.DatabasePassword+"@tcp("+info.DatabaseAddress+")/?parseTime=True&loc=Local&charset=utf8mb4")
	if err != nil{
		if strings.Contains(err.Error(),"Unknown character set"){
			//for version lower 5.5
			db, err = sql.Open("mysql", info.DatabaseUser+":"+info.DatabasePassword+"@tcp("+info.DatabaseAddress+")/?parseTime=True&loc=Local&charset=utf8")
			if err!=nil{
				log.Fatalf("err: %v\n", err)
			}
		}else{
			log.Fatalf("err: %v\n", err)
		}
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	color.Infoln(info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress," connection inited successfully.")
	return db
}

/*
  Mssql database initialization, and return database handle
  @Param  info (the information user inputs)
  @Return sql.DB (database handle)
*/
func MssqlDBinit(info *models.InitData) (*sql.DB) {
	var db *sql.DB
	if info.WindowsAuth{
		port:="1433"
		server:="127.0.0.1"

		if info.DatabaseAddress!=""{
			val := strings.Split(info.DatabaseAddress, ":")
			server=val[0]
			port=val[1]
		}
		connString := fmt.Sprintf("server=%s;port=%s;encrypt=disable;trusted_connection=yes;",server,port)
		dbF, err := sql.Open("mssql", connString)
		if err != nil {
			log.Printf("Error connecting mssql:  %s:%s\n", server,port)
			log.Fatalf(err.Error())
		}
		db=dbF
	}else{
		if info.DatabaseAddress==""{
			log.Fatalf("Enter the database address,please!")
		}

		dbF, err := sql.Open("mssql", fmt.Sprintf("sqlserver://%v:%v@%v/?connection&encrypt=disable&charset=utf8mb4", info.DatabaseUser,url.PathEscape(info.DatabasePassword),info.DatabaseAddress))

		if err != nil{
			log.Fatalf("err: %v\n", err)
		}
		db=dbF
	}
	err := db.Ping()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	color.Infoln(info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress," connection inited successfully.")
	return db
}

/*
  Oracle database initialization, and return database handle
  @Param  info (the information user inputs)
  @Return sql.DB (database handle)
*/
func OracleDBinit(info *models.InitData) (*sql.DB) {
	var db *sql.DB
	var err error
	if info.DatabaseAddress==""&&info.TNSFile==""{
		log.Fatalf("Enter the database address or tnsnames.ora config file path, please!")
	}

	var connString  string
	index := strings.Index(info.DatabaseAddress, "/")
	if index != -1 {
		connString=fmt.Sprintf("oracle://"+info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress+info.DatabaseInstance)
		} else {
			if info.DatabaseInstance!=""{
				connString=fmt.Sprintf("oracle://"+info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress+"/"+info.DatabaseInstance)
			}else{
				info.DatabaseInstance = "/ORCL"
				connString=fmt.Sprintf("oracle://"+info.DatabaseUser+":"+info.DatabasePassword+"@"+info.DatabaseAddress+info.DatabaseInstance)
			}
	}

	db, err = sql.Open("oracle", connString)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	connectInfo:=strings.Split(connString,"://")[1]
	color.Infoln(connectInfo," connection inited successfully.")

	return db
}

// check the connectivity with oracle database setting 5 seconds timeout
func CheckConnectivity(connStr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		return false
	}
	return true
}

/*
  mysql: count all tables and add them into list
  @Param  db (database handle)
  @Return []string (all the tables in the database)
*/
func CountAllTables(db *sql.DB) ([]string){
	var tableList []string
	//Count Number of Rows in each table
	rows, err := db.Query("select table_schema,table_name from information_schema.tables;")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rows.Close()
	cols,err:=rows.Columns()
	if err!=nil{
		log.Fatalf(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next(){
		// read the row on the table
		// each column value will be stored in the slice
		err = rows.Scan(scanArgs...)
		utils.CheckError("Error scanning rows from table", err)

		if string(values[0])=="information_schema"||string(values[0])=="mysql"||string(values[0])=="performance_schema"||string(values[0])=="sys"{
			continue
		}
		tableList=append(tableList,string(values[0])+"."+string(values[1]))
	}
	return tableList
}

/*
  mssql: count all tables and add them into list
  @Param  db (database handle)
  @Return []string (all the tables in the database)
*/
func CountAllTablesMs(db *sql.DB) ([]string){
	var tableList []string
	dbs:=QueryWrapped(db,"mssql","database","","",0)
	for dbs.Next() {
		var dbsname string
		err := dbs.Scan(&dbsname)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if dbsname=="tempdb"||dbsname=="master"||dbsname=="model"||dbsname=="msdb"||dbsname=="ReportServer"||dbsname=="ReportServerTempDB"{
			continue
		}

		tables:= QueryWrapped(db,"mssql","table",dbsname,"",0)
		var tblname string
		for tables.Next() {
			err = tables.Scan(&tblname)
			if err != nil {
				log.Fatalf(err.Error())
			}
			tableList=append(tableList,dbsname+"."+tblname)
		}
	}
	return tableList
}

/*
  oracle: count all tables and add them into list
  @Param  db (database handle)
  @Return []string (all the tables in the database)
*/
func CountAllTablesOra(db *sql.DB) ([]string){
	var tableList []string
	//list all the tables in the database
	rows, err := db.Query("SELECT OWNER,TABLE_NAME FROM all_tables ORDER BY OWNER")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rows.Close()
	cols,err:=rows.Columns()
	if err!=nil{
		log.Fatalf(err.Error())
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next(){
		// read the row on the table
		// each column value will be stored in the slice
		err = rows.Scan(scanArgs...)
		utils.CheckError("Error scanning rows from table", err)

		if string(values[0])=="XDB"||string(values[0])=="CTXSYS"||string(values[0])=="DBSNMP"||string(values[0])=="EXFSYS"||string(values[0])=="MDSYS"||string(values[0])=="ORDSYS"||string(values[0])=="OLAPSYS"||string(values[0])=="SYS"||string(values[0])=="SYSMAN"||string(values[0])=="FLOWS_FILES"||string(values[0])=="APEX_030200"||string(values[0])=="APPQOSSYS"||string(values[0])=="PM"||string(values[0])=="ORDDATA"||string(values[0])=="IX"||string(values[0])=="WMSYS"||string(values[0])=="OWBSYS"||string(values[0])=="OUTLN"||string(values[0])=="SYSTEM"||string(values[0])=="OE"||string(values[0])=="SH"||string(values[0])=="HR"{
			continue
		}
		tableList=append(tableList,string(values[0])+"."+string(values[1]))
	}
	return tableList
}

/*
  db.Query wrapped function
  @Param  db (database handle)
  @Param  typeD (the type of database)
  @Param  queryType (the query type passes into database query wrapped function)
  @Param  dbsname (the database name or the schema name)
  @Param  tblname (the table name)
  @Param  num (the number of rows returned from database)
  @Return *sql.Rows (the information from database)
*/
func QueryWrapped(db *sql.DB,typeD string,queryType string,dbsname string,tblname string,num int) *sql.Rows{
	switch queryType{
	case "database":
		if typeD=="mysql"{
			return QueryFromDatabases(db,"SHOW DATABASES")
		}else if typeD=="mssql"{
			return QueryFromDatabases(db,"SELECT NAME FROM MASTER.DBO.SYSDATABASES")
		}else if typeD=="oracle"{
			return QueryFromDatabases(db,"SELECT USERNAME FROM ALL_USERS")
		}
	case "table":
		if typeD=="mysql"{
			return QueryFromDatabases(db,"SElECT TABLE_NAME from information_schema.tables where table_schema = '"+dbsname+"'")
		}else if typeD=="mssql"{
			return QueryFromDatabases(db,"SELECT NAME FROM " + dbsname + ".sys.tables")
		}else if typeD=="oracle"{
			return QueryFromDatabases(db,"SELECT TABLE_NAME FROM ALL_TABLES WHERE OWNER='"+dbsname+"'")
		}
	case "column":
		if typeD=="mysql"{
			return QueryFromDatabases(db,"SELECT COLUMN_NAME  FROM information_schema.columns WHERE table_name='"+tblname+"'and table_schema='"+dbsname+"'")
		}else if typeD=="mssql"{
			return QueryFromDatabases(db,"SELECT COLUMN_NAME  FROM "+dbsname+".information_schema.columns where table_name = '"+tblname+"'")
		}else if typeD=="oracle"{
			return QueryFromDatabases(db," SELECT t1.COLUMN_NAME FROM all_tab_columns t1 WHERE OWNER='"+dbsname+"' AND table_name='"+tblname+"'")
		}
	case "data":
		if typeD=="mysql"{
			return QueryFromDatabases(db,"SELECT * FROM " + dbsname+"."+fmt.Sprintf("`%s`", tblname)+" LIMIT " + strconv.Itoa(num))
		}else if typeD=="mssql"{
			return QueryFromDatabases(db,"SELECT TOP "+strconv.Itoa(num)+" * FROM "+dbsname+".."+tblname)
		}else if typeD=="oracle"{
			return QueryFromDatabases(db,"SELECT t1.* FROM " + dbsname+"."+tblname +" t1 WHERE ROWNUM <="+strconv.Itoa(num))
		}
	}
	return nil
}

/*
  db.QueryRow wrapped function
  @Param  sql.DB (database handle)
  @Param  typeD (the type of database)
  @Param  dbsname (the database name or the schema name)
  @Param  tblname (the table name)
  @Param  *count (get the result of database's count function)
*/
func QueryCount(db *sql.DB,typeD string,dbsname string,tblname string,count *int){
	if typeD=="mysql"{
		db.QueryRow("select count(*) from " +dbsname+"."+tblname).Scan(count)
	}else if typeD=="mssql"{
		db.QueryRow("select count(*) from "+dbsname+".."+tblname).Scan(count)
	}else if typeD=="oracle"{
		db.QueryRow("select count(*) from " +dbsname+"."+tblname).Scan(count)
	}
}

/*
  db.Query function
  @Param  sql.DB (database handle)
  @Param  queryStr (the string of querying the database)
  @Return *sql.Rows (the information from database)
*/
func QueryFromDatabases(db *sql.DB,queryStr string) *sql.Rows{
	rows,err := db.Query(queryStr)
	if err != nil {
		fmt.Println("Err query: ",queryStr)
		log.Fatal(err.Error())
	}
	return  rows
}