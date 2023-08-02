package utils

import (
	"fmt"
	"log"
	"strings"
)

//For the name of final output file
type InfoStruct struct {
	IP string
	Port string
	User string
}

/*
  Deal with the output file name
  @Param  address (database address)
  @Param  user (database user)
  @Param  datatype (database type)
  @Return InforStruct
*/
func OutputFileName(address,user,datatype string) InfoStruct{
	defer func(){
		r := recover()
		if r != nil {
			fmt.Println("Please verify your input!")
			log.Fatal("PANIC :", r)
		}
	}()
	if address==""||user==""{
		if datatype=="mongo"{
			user="NULL"
		}else{
			log.Fatalf("Please input the database address and database user!")
		}
	}
	res := strings.Split(address, ":")
	ip := res[0]
	port := res[1]
	return InfoStruct{IP:ip,Port: port,User: user}
}


// helper function with message to handle errors
func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
// helper function to handle errors
func CheckErrorExit(err error){
	if err != nil {
		log.Fatalf(err.Error())
	}
}