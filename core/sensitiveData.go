package core

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"database/sql"
	"dataMiner/models"
	"dataMiner/utils"
	"dataMiner/dblib"
	"github.com/dlclark/regexp2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Turn tables into iterator struct
type  Alltables struct {
	tables []string
}

func (i Alltables) Iterator() *Iterator {
	return &Iterator{
		data:  i,
		index: 0,
	}
}

type Iterator struct {
	data Alltables
	index int
}

func (i *Iterator) HasNext() bool {
	return i.index < len(i.data.tables)
}

func (i *Iterator) Next() string {
	tbl := i.data.tables[i.index]
	i.index++
	return tbl
}

var passwdTable []string   //store the table which has pass,passwd or password column

/*
  Look for sensitive data in the database,such as email address,phone number,ID card number and password
  @Param  db (database handle)
  @Param  client (mongo database handle)
  @Param  connectionString (postgre database connection string)
  @Param  tableList (all the tables in the database)
  @Param  num (the number of rows returned from database)
  @Param  thread (the thread user specified)
  @Param  outputID (the output file name)
  @Param  pattern (the regular expresstion pattern)
  @Param  typeD (the type of database)
*/
func LookforSensitiveData(db *sql.DB,client *mongo.Client,connectionString string,tableList []string,num int, thread int,outputID utils.InfoStruct,pattern string,typeD string){
	fmt.Println("Task is in processing......")
	var results []models.SensitiveData  //store the sensitive data which matches the pattern

	if typeD=="mysql"{
		Searchpasswd(db)
	}else if typeD=="mssql"{
		SearchpasswdMs(db)
	}else if typeD=="oracle"{
		SearchpasswdOra(db)
	}else if typeD=="postgre"{
		SearchpasswdPs(connectionString,tableList)
	}

	tables := Alltables{tables: tableList}
	it := tables.Iterator()

	wg := &sync.WaitGroup{}
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		if typeD=="mongo"{
			go StartSearchDataMongo(it,wg,client,num,pattern,typeD,&results)
		}else if typeD=="postgre"{
			go StartSearchDataPs(it,wg,connectionString,num,pattern,typeD,&results)
		}else{
			go StartSearchData(it,wg,db,num,pattern,typeD,&results)
		}
	}
	wg.Wait()
	Output(outputID,results)
}

/*
  Start searching sensitive data using goroutines
  @Param  it (the iterator of table list)
  @Param  wg (the WaitGroup of goroutines)
  @Param  db (database handle)
  @Param  num (the number of rows returned from database)
  @Param  pattern (the regular expresstion pattern)
  @Param  typeD (the type of database)
  @Param  results (the result of searching)
*/
func StartSearchData(it *Iterator,wg *sync.WaitGroup,db *sql.DB,num int,pattern string,typeD string,results *[]models.SensitiveData) {
	defer wg.Done()
	for{
		if(it.HasNext()) {
			tbl:=it.Next()
			// Get data from each table
			val := strings.Split(tbl, ".")
			data:=dblib.QueryWrapped(db ,typeD,"data",val[0],val[1],num)
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

			for data.Next() {
				// read the row on the table
				// each column value will be stored in the slice
				err = data.Scan(scanArgs...)
				utils.CheckError("Error scanning rows from table", err)

				var value string
				for _, col := range values {
					// Here  check if the value is nil (NULL value)
					if col == nil {
						continue
					} else {
						value = string(col)
						var found string
						if pattern!=""{
							found=MatchSensitivityUserDefined(value,pattern)
						}else{
							found = MatchSensitivity(value,tbl,typeD)
						}
						if found != "" {
							//store it in results
							val:=strings.Split(tbl,".")
							rtmp:=models.SensitiveData{DatabaseName: val[0],TableName: val[1],Data: value,Type: found}
							*results=append(*results,rtmp)
						}
					}
				}
			}
		}else{
			break
		}
	}
}

/*
  Start searching sensitive data in Mongodb database using goroutines
  @Param  it (the iterator of table list)
  @Param  wg (the WaitGroup of goroutines)
  @Param  client (database handle)
  @Param  num (the number of rows returned from database)
  @Param  pattern (the regular expresstion pattern)
  @Param  typeD (the type of database)
  @Param  results (the result of searching)
*/
func StartSearchDataMongo(it *Iterator,wg *sync.WaitGroup,client *mongo.Client,num int,pattern string,typeD string,results *[]models.SensitiveData) {
	defer wg.Done()
	for {
		if(it.HasNext()) {
			//get collection from iterator
			tbl := it.Next()
			val := strings.Split(tbl, ".")
			var resultsMongo []bson.M
			dblib.GetDocuments(client,val[0],val[1],num,&resultsMongo)

			//get the documents from each collection
			for _, result := range resultsMongo {
				iter := reflect.ValueOf(result).MapRange()
				for iter.Next() {
					key := iter.Key().String()
					value := iter.Value().Interface()
					var keyToCheck,valueToCheck string
					if value!=""{
						keyToCheck,valueToCheck=SingleMongoData(key+".",value)
					}
					var found string
					if pattern!=""{
						found=MatchSensitivityUserDefined(valueToCheck,pattern)
					}else{
						found = MatchSensitivity(valueToCheck,keyToCheck,typeD)
					}
					if found != "" {
						//store it in results
						val:=strings.Split(tbl,".")
						rtmp:=models.SensitiveData{DatabaseName: val[0],TableName: val[1],Data: keyToCheck+" : "+valueToCheck,Type: found}
						*results=append(*results,rtmp)
					}
				}

			}
		}else{
			break
		}
	}
}

/*
  Start searching sensitive data in Postgre database using goroutines
  @Param  it (the iterator of table list)
  @Param  wg (the WaitGroup of goroutines)
  @Param  connectionString (postgre connection string)
  @Param  num (the number of rows returned from database)
  @Param  pattern (the regular expresstion pattern)
  @Param  typeD (the type of database)
  @Param  results (the result of searching)
*/
func StartSearchDataPs(it *Iterator,wg *sync.WaitGroup,connectionString string,num int,pattern string,typeD string,results *[]models.SensitiveData){
	defer wg.Done()
	for{
		if(it.HasNext()) {
			tbl:=it.Next()

			parts := strings.Split(tbl, ".")
			database := parts[0]
			schema := parts[1]
			table := strings.Join(parts[2:], ".")

			// Get data from each table
			data:=dblib.QueryDataPs(connectionString,database,schema,table,num)
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

			for data.Next() {
				// read the row on the table
				// each column value will be stored in the slice
				err = data.Scan(scanArgs...)
				utils.CheckError("Error scanning rows from table", err)

				var value string
				for _, col := range values {
					// Here  check if the value is nil (NULL value)
					if col == nil {
						continue
					} else {
						value = string(col)
						var found string
						if pattern!=""{
							found=MatchSensitivityUserDefined(value,pattern)
						}else{
							found = MatchSensitivity(value,tbl,typeD)
						}
						if found != "" {
							//store it in results
							rtmp:=models.SensitiveData{DatabaseName: database,TableName: schema+"."+table,Data: value,Type: found}
							*results=append(*results,rtmp)
						}
					}
				}
			}
		}else{
			break
		}
	}
}

/*
  Matching sensitive data entrance
  @Param  data (the data to be verified)
  @Param  dtname (the table to be verified, check whether it's in the table list whose column name contains password, passwd or pass)
  @Return sting type (the type of sensitive data)
*/
func MatchSensitivity(data, dtname, typeD string) string{
	if VerifyEmailFormat(data){
		return "Email"
	}else if VerifyMobileFormat(data){
		return "Mobile"
	}else if VerifyIDFormat(data){
		return "IDCardNumer"
	}else if isInPasswdSlice(dtname,typeD){
		if VerifyPasswdFormat(data){
			return "Password"
		}else{
			return ""
		}
	}
	return ""
}

/*
  Matching user-defined data
  @Param  data (the data to be verified)
  @Return sting type (the type of sensitive data)
*/
func MatchSensitivityUserDefined(data string,pattern string) string{
	if VerifyPatternFormat(data,pattern){
		return "UserDefined"
	}
	return ""
}

/*
  Check whether the table is in the password table list or whether the key is related to the password field
  @Param  dtname (the table to be verified, check whether it's in the table list whose column name contains password,passwd or pass, or the key in mongodb whose name contains password,passwd or pass
  @Param  typeD (the type of database)
  @Return bool type
*/
func isInPasswdSlice(dtname, typeD string) bool{
	if typeD=="mongo"{
		if strings.Contains(strings.ToLower(dtname),"passwd")||strings.Contains(strings.ToLower(dtname),"password")||strings.Contains(strings.ToLower(dtname),"pass"){
			return true
		}
	}else{
		for _, s := range passwdTable{
			if s == dtname {
				return true
			}
		}
	}
	return false
}

/*
  Save the result of searching sensitive data
  @Param  outputID (the output file name)
  @Param  results (the result of searching)
*/
func Output(outputID utils.InfoStruct,results []models.SensitiveData){
	var csv []models.SensitiveData
	for _,i:=range results {
		tmp:=models.SensitiveData{DatabaseName: i.DatabaseName,TableName: i.TableName,Data: i.Data,Type: i.Type}
		csv=append(csv,tmp)
	}
	//save the result to csv
	utils.SavetocsvD(csv,outputID)
	//save the result to html
	utils.SavetohtmlD(csv,outputID)
}

/*
  Matching email like data
  @Param  email (the data to be verified)
  @Return bool type
*/
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

/*
  Matching phone number like data
  @Param  mobileNum (the data to be verified)
  @Return bool type
*/
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

/*
  Matching ID card number like data
  @Param  idNum (the data to be verified)
  @Return bool type
*/
func VerifyIDFormat (idNum string) bool{
	pattern:= `(^[1-9]\d{5}[1-9]\d{3}(((0[2])([0|1|2][0-8])|(([0-1][1|4|6|9])([0|1|2][0-9]|[3][0]))|(((0[1|3|5|7|8])|(1[0|2]))(([0|1|2]\d)|3[0-1]))))((\d{4})|\d{3}[Xx])$)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(idNum)
}

//username verify     high false positive, not add currently
func VerifyUsernameFormat (name string) bool{
	pattern:=`^[\x{4e00}-\x{9fa5}]{2,4}$`
	pattern1:=`^[a-zA-Z0-9_-]{3,15}$`
	reg := regexp.MustCompile(pattern)
	reg1:=regexp.MustCompile(pattern1)
	return reg.MatchString(name)|| reg1.MatchString(name)
}

/*
  Matching password like data
  @Param  passwd (the data to be verified)
  @Return bool type
*/
func VerifyPasswdFormat (passwd string) bool{
	expr := `(?![A-Z]+$)(?![a-z]+$)(?![0-9]+$)(?![-!@#$^&+,.]+$)[\w-!@#$^&+,.]{5,32}$`
	expr1:=`[a-zA-Z0-9]{5,}`
	re := regexp2.MustCompile(expr, 0)
	re1:=regexp2.MustCompile(expr1,0)
	isMatch, _ := re.MatchString(passwd)
	isMatch1, _ := re1.MatchString(passwd)
	if isMatch||isMatch1{
		return true
	}
	return false
}

/*
  User defined regular expression func
  @Param  data (the data to be verified)
  @Param  pattern (the regular expresstion pattern)
  @Return bool type
*/
func VerifyPatternFormat(data string,pattern string) bool{
	re := regexp2.MustCompile(pattern, 0)
	isMatch, _ := re.MatchString(data)
	if isMatch{
		return true
	}
	return false
}

/*
  mysql database: search the table which has password column and store it in passwdTable slice
  @Param  db (database handle)
*/
func Searchpasswd(db *sql.DB){

	checkList:=[]string{"passwd","password","pass"}
	//get all the column names in the databases
	rows, err := db.Query("SELECT table_schema,table_name,column_name FROM information_schema.columns")
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

		for _,key:=range checkList {
			if strings.Contains(strings.ToLower(string(values[2])),key){
				passwdTable=append(passwdTable, string(values[0])+"."+string(values[1]))
			}
		}
	}
}

/*
  mssql database: search the table which has password column and store it in passwdTable slice
  @Param  db (database handle)
*/
func SearchpasswdMs(db *sql.DB){

	checkList:=[]string{"passwd","password","pass"}

	dbs:= dblib.QueryWrapped(db,"mssql","database","","",0)  //list all the databases
	for dbs.Next(){
		var dbsname string
		err := dbs.Scan(&dbsname)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if dbsname=="tempdb"||dbsname=="master"||dbsname=="model"||dbsname=="msdb"||dbsname=="ReportServer"||dbsname=="ReportServerTempDB"{
			continue
		}

		tables:= dblib.QueryWrapped(db,"mssql","table",dbsname,"",0)
		var tblname string
		for tables.Next() {
			err = tables.Scan(&tblname)
			if err != nil {
				log.Fatalf(err.Error())
			}

			columnRows:=dblib.QueryWrapped(db ,"mssql","column",dbsname,tblname,0)
			defer columnRows.Close()

			for columnRows.Next() {
				var columnName string
				if err := columnRows.Scan(&columnName); err != nil { log.Fatal(err) }
				for _,key:=range checkList {
					if strings.Contains(strings.ToLower(columnName),key){
						passwdTable=append(passwdTable, dbsname+"."+tblname)
					}
				}
			}
		}
	}
}

/*
  oracle database: search the table which has password column and store it in passwdTable slice
  @Param  db (database handle)
*/
func SearchpasswdOra(db *sql.DB){
	checkList:=[]string{"passwd","password","pass"}
	//get all the column names in the databases
	rows, err := db.Query("SELECT OWNER,TABLE_NAME,COLUMN_NAME FROM all_tab_columns")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	cols,err:=rows.Columns()
	if err!=nil{
		panic(err.Error())
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

		for _,key:=range checkList {
			if strings.Contains(strings.ToLower(string(values[2])),key){
				passwdTable=append(passwdTable, string(values[0])+"."+string(values[1]))
			}
		}
	}
}

/*
  postgre database: search the table which has password column and store it in passwdTable slice
  @Param  connectionString (postgre connection string)
  @Param  tableList (all the tables in the database)
*/
func SearchpasswdPs(connectionString string,tableList []string){
	checkList:=[]string{"passwd","password","pass"}
	for _,tbl:=range tableList{
		parts := strings.Split(tbl, ".")
		database := parts[0]
		schema := parts[1]
		table := strings.Join(parts[2:], ".")

		// get column names
		columnRows:=dblib.QueryColumnsPs(connectionString,database,schema,table)
		defer columnRows.Close()

		//Loop through the rows and append the column names to the columnNames
		for columnRows.Next() {
			var columnName string
			if err := columnRows.Scan(&columnName); err != nil {log.Fatal(err) }
			for _,key:=range checkList {
				if strings.Contains(strings.ToLower(columnName),key){
					passwdTable=append(passwdTable, database+"."+schema+"."+table)
				}
			}
		}
	}
}

/*
  Get the key and data from one document
  @Param  prefix (the key of document)
  @Param  doc (the value of document)
  @Return  key and value (return the final key and value of document)
*/
func SingleMongoData(prefix string, doc interface{}) (string,string) {
	// Get the type of the value
	value := reflect.ValueOf(doc)
	var valueReturn string
	switch value.Kind() {
	case reflect.Map:
		// Iterate over each key-value pair in the map
		iter := reflect.ValueOf(doc).MapRange()
		for iter.Next() {
			key := iter.Key().String()
			value := iter.Value().Interface()
			// Recursively get the keys in the value
			SingleMongoData(prefix+key+".", value)
		}
	case reflect.Slice:
		// Iterate over each element in the slice
		for i := 0; i < value.Len(); i++ {
			// Recursively get the keys in the element
			SingleMongoData(prefix+fmt.Sprintf("[%d]", i)+".", value.Index(i).Interface())
		}
	default:
		//get the final key and value
		prefix=strings.Trim(prefix,".")
		valueReturn=fmt.Sprint(value)
	}
	return prefix,valueReturn
}