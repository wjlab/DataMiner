package core

import (
	"bufio"
	"log"
	"os"
	"strings"
)
type InputInfo struct{
	Schema string
	Address string
	User string
	Passwd string
	AuthSource string     //Mongodb authentication need to provie database name
}

/*
  Read all databases connection information from file and store them in slice
  @Param  filename (the file where the databases connection information stored)
  @Return []InputInfo struct (the struct array passes into next function)
*/
func Batch(filename string) []InputInfo{
	var inputs []InputInfo
	databases,err:=Readfile(filename)
	if err!=nil{
		log.Fatal(err)
	}
	for _,i:=range databases{
		tmp:=SplitInfo(i)
		inputs=append(inputs,tmp)
	}
	return inputs

}

/*
  Split the information from file into formal struct
  @Param  str (the string of each row in the file)
  @Return InputInfo struct (the struct passes into next function)
*/
func SplitInfo(str string) InputInfo{
	var authSource,user,password string
	strs := strings.Split(str, "://")
	schema := strs[0]

	lastIndex := strings.LastIndex(strs[1], "@")
	address:=strs[1][lastIndex+1:]
	userinfo:=strs[1][:lastIndex]

	if strings.Contains(userinfo,":"){
		res := strings.Split(userinfo, ":")
		user = res[0]
		password = res[1]
	}

	if strings.Contains(address,"?"){
		value:=strings.Split(address,"?")

		if len(value)==2{
			address=strings.Trim(value[0],"/")
			authSource=value[1]
		}else{
			log.Fatal("Please provide right mongodb format when mongodb need database name to authenticate, like: 127.0.0.1:27017?databaseName")
		}
	}

	return InputInfo{Schema: schema,User: user,Passwd: password,Address: address,AuthSource: authSource}
}

/*
  Read the batch databases information from file
  @Param  filename (the file where the databases connection information stored)
  @Return []string (the information of each databases in the file)
  @Return error
*/
func Readfile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewScanner(f)
	var result []string

	for {
		if !buf.Scan() {
			break
		}

		line := buf.Text()
		line = strings.TrimSpace(line)

		result = append(result, line)
	}
	return result, nil
}