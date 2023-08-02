package core

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"fmt"
	"os/user"
	"database/sql"
	"dataMiner/models"
	"github.com/dlclark/regexp2"
	"github.com/urfave/cli"
	"golang.org/x/net/proxy"
	"dataMiner/utils"
)


var num int          //the number for extract information from databases
var thread int       //the number of thread for SearchSensitiveData function
var pattern string   //the user-defined regular expression for SearchSensitiveData function
var databaseType    string   //the type of database
var databaseAddress string
var databaseUser    string
var databasePassword  string
var singleTable       string
var proxyAddress      string
var proxyUser         string
var proxyPassword     string
var filename          string
var windowsAuth       bool
var proxyConnection    proxy.Dialer  //socks5 proxy function needs
func Execute() {
	Logo()
	app := &cli.App{
		Name: "DataMiner",
		Usage: "The tool used to extract the information from databases quickly.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "databaseType",
				Aliases: []string{"T"},
				Usage:   "-T mysql (currently supports mysql,mssql,oracle and mongo)",
				Destination: &databaseType,
			},
			&cli.StringFlag{
				Name:    "databaseAddress",
				Aliases: []string{"da"},
				Usage:   "-da 127.0.0.1:3306",
				Destination: &databaseAddress,
			},
			&cli.StringFlag{
				Name:    "databaseUser",
				Aliases: []string{"du"},
				Usage:   "-du name",
				Destination: &databaseUser,
			},
			&cli.StringFlag{
				Name:    "databasePassword",
				Aliases: []string{"dp"},
				Usage:   "-dp passwd",
				Destination: &databasePassword,
			},
			&cli.StringFlag{
				Name:    "proxyAddress",
				Aliases: []string{"pa"},
				Usage:   "-pa 127.0.0.1:8080",
				Destination: &proxyAddress,
			},
			&cli.StringFlag{
				Name:    "proxyUser",
				Aliases: []string{"pu"},
				Usage:   "-pu name",
				Destination: &proxyUser,
			},
			&cli.StringFlag{
				Name:    "proxyPassword",
				Aliases: []string{"pp"},
				Usage:   "-pp passwd",
				Destination: &proxyPassword,
			},
			&cli.StringFlag{
				Name:    "databaseTable",
				Aliases: []string{"dt"},
				Usage:   "-dt database.table",
				Destination: &singleTable,
			},
			&cli.StringFlag{
				Name:    "fileInput",
				Aliases: []string{"f"},
				Usage:   "-f filename(like: -f test.txt)",
				Destination: &filename,
			},
			&cli.BoolFlag{
				Name:    "WindowsAuth",
				Aliases: []string{"WA"},
				Usage:   "-WA (Only for mssql, if choose this, it will connect mssql using windows authentication)",
				Value: false,
				Destination: &windowsAuth,
			},

			&cli.IntFlag{
				Name:    "num",
				Aliases: []string{"n"},
				Value:   3,
				Usage:   "-n (The number for extract information from databases)",
				Destination: &num,
			},
			&cli.IntFlag{
				Name:    "thread",
				Aliases: []string{"t"},
				Value:   5,
				Usage:   "-t 1 (Only For SearchSensitiveData function)",
				Destination: &thread,
			},
			&cli.StringFlag{
				Name:    "pattern",
				Aliases: []string{"p"},
				Usage:   "-p pattern(Only For SearchSensitiveData function,like searching for username: -p ^[\\x{4e00}-\\x{9fa5}]{2,4}$ )",
				Destination: &pattern,
			},

		},
		Commands: []*cli.Command{
			{
				Name:    "Sampledata",
				Aliases: []string{"SD"},
				Usage:   "Command for getting Samledata from databases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "databaseType",
						Aliases: []string{"T"},
						Usage:   "-T mysql (currently supports mysql,mssql,oracle and mongo)",
						Destination: &databaseType,
					},
					&cli.StringFlag{
						Name:    "databaseAddress",
						Aliases: []string{"da"},
						Usage:   "-da 127.0.0.1:3306)",
						Destination: &databaseAddress,
					},
					&cli.StringFlag{
						Name:    "databaseUser",
						Aliases: []string{"du"},
						Usage:   "-du name",
						Destination: &databaseUser,
					},
					&cli.StringFlag{
						Name:    "databasePassword",
						Aliases: []string{"dp"},
						Usage:   "-dp passwd",
						Destination: &databasePassword,
					},
					&cli.StringFlag{
						Name:    "proxyAddress",
						Aliases: []string{"pa"},
						Usage:   "-pa 127.0.0.1:8080",
						Destination: &proxyAddress,
					},
					&cli.StringFlag{
						Name:    "proxyUser",
						Aliases: []string{"pu"},
						Usage:   "-pu name",
						Destination: &proxyUser,
					},
					&cli.StringFlag{
						Name:    "proxyPassword",
						Aliases: []string{"pp"},
						Usage:   "-pp passwd",
						Destination: &proxyPassword,
					},
					&cli.StringFlag{
						Name:    "fileInput",
						Aliases: []string{"f"},
						Usage:   "-f filename(like: -f test.txt)",
						Destination: &filename,
					},
					&cli.BoolFlag{
						Name:    "WindowsAuth",
						Aliases: []string{"WA"},
						Usage:   "-WA (Only for mssql, if choose this, it will connect mssql using windows authentication)",
						Value: false,
						Destination: &windowsAuth,
					},
					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Value:   3,
						Usage:   "-n (The number for extract information from databases)",
						Destination: &num,
					},
				},
				Action:  func(c *cli.Context) error {
					var start = time.Now()
					iniInfo:=initData()
					//establish proxy connection
					var connection  net.Conn
					if proxyAddress!=""&&filename==""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
						connection= ProxyConnect(databaseAddress)
					}else if proxyAddress!=""&&filename!=""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
					}
					defer func(){
						if connection != nil {
							_ = connection.Close()
						}
					}()

					if windowsAuth{
						if databaseAddress==""{
							databaseAddress="127.0.0.1:1433"
						}
						user,_ := user.Current()
						userName:=strings.Split(user.Username, "\\")
						databaseUser=userName[len(userName)-1]
					}

					if filename==""{
						outputID:=utils.OutputFileName(databaseAddress,databaseUser,databaseType)
						SingleSampleData(outputID,iniInfo)
                        }else{
						inputs:=Batch(filename)

						for n,j:=range inputs {
							var connectionTmp net.Conn
							iniInfoB:=models.InitData{DatabaseType: j.Schema,DatabaseAddress: j.Address,DatabaseUser: j.User,DatabasePassword: j.Passwd,AuthSource: j.AuthSource}
							fmt.Print("No." + strconv.Itoa(n+1) + " ")
							OutputID := utils.OutputFileName( j.Address, j.User,j.Schema)
							//establish proxy connection
							if proxyAddress!=""{
								connectionTmp=ProxyConnect(j.Address)
							}
							SingleSampleData(OutputID,iniInfoB)
							if connectionTmp!=nil{
								err:=connectionTmp.Close()
								if err!=nil{
									log.Fatal(err)
								}
							}

						}
					}
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end)
					return nil
				},
			},
			{
				Name:    "Overview",
				Aliases: []string{"OV"},
				Usage:   "Command for overviewing the databases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "databaseType",
						Aliases: []string{"T"},
						Usage:   "-T mysql (currently supports mysql,mssql,oracle and mongo)",
						Destination: &databaseType,
					},
					&cli.StringFlag{
						Name:    "databaseAddress",
						Aliases: []string{"da"},
						Usage:   "-da 127.0.0.1:3306",
						Destination: &databaseAddress,
					},
					&cli.StringFlag{
						Name:    "databaseUser",
						Aliases: []string{"du"},
						Usage:   "-du name",
						Destination: &databaseUser,
					},
					&cli.StringFlag{
						Name:    "databasePassword",
						Aliases: []string{"dp"},
						Usage:   "-dp passwd",
						Destination: &databasePassword,
					},
					&cli.StringFlag{
						Name:    "proxyAddress",
						Aliases: []string{"pa"},
						Usage:   "-pa 127.0.0.1:8080",
						Destination: &proxyAddress,
					},
					&cli.StringFlag{
						Name:    "proxyUser",
						Aliases: []string{"pu"},
						Usage:   "-pu name",
						Destination: &proxyUser,
					},
					&cli.StringFlag{
						Name:    "proxyPassword",
						Aliases: []string{"pp"},
						Usage:   "-pp passwd",
						Destination: &proxyPassword,
					},
					&cli.StringFlag{
						Name:    "fileInput",
						Aliases: []string{"f"},
						Usage:   "-f filename(like: -f test.txt)",
						Destination: &filename,
					},
					&cli.BoolFlag{
						Name:    "WindowsAuth",
						Aliases: []string{"WA"},
						Usage:   "-WA (Only for mssql, if choose this, it will connect mssql using windows authentication)",
						Value: false,
						Destination: &windowsAuth,
					},

				},
				Action:  func(c *cli.Context) error {
					var start = time.Now()
					iniInfo:=initData()
					//establish proxy connection
					var connection  net.Conn
					if proxyAddress!=""&&filename==""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
						connection= ProxyConnect(databaseAddress)
					}else if proxyAddress!=""&&filename!=""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
					}
					defer func(){
						if connection != nil {
							_ = connection.Close()
						}
					}()

					if windowsAuth{
						if databaseAddress==""{
							databaseAddress="127.0.0.1:1433"
						}
						user,_ := user.Current()
						userName:=strings.Split(user.Username, "\\")
						databaseUser=userName[len(userName)-1]
					}

					if filename==""{
						outputID:=utils.OutputFileName(databaseAddress,databaseUser,databaseType)
						SingleOverview(outputID,iniInfo)
					}else{
						inputs:=Batch(filename)
						for n,j:=range inputs{
							var connectionTmp net.Conn
							databaseUser=j.User
							databasePassword=j.Passwd
							databaseAddress=j.Address
							iniInfoB:=models.InitData{DatabaseType: j.Schema,DatabaseAddress: j.Address,DatabaseUser: j.User,DatabasePassword: j.Passwd,AuthSource: j.AuthSource}
							fmt.Print("No."+strconv.Itoa(n+1)+" ")
							outputID:=utils.OutputFileName(j.Address,j.User,j.Schema)
							//establish proxy connection
							if proxyAddress!="" {
								connectionTmp = ProxyConnect(j.Address)
							}
							SingleOverview(outputID,iniInfoB)
							if connectionTmp!=nil{
								err:=connectionTmp.Close()
								if err!=nil{
									log.Fatal(err)
								}
							}
						}
					}
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end)
					return nil
				},
			},

			{
				Name:    "SearchSensitiveData",
				Aliases: []string{"SS"},
				Usage:   "Command for searching sensitive data from databases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "databaseType",
						Aliases: []string{"T"},
						Usage:   "-T mysql (currently supports mysql,mssql,oracle and mongo)",
						Destination: &databaseType,
					},
					&cli.StringFlag{
						Name:    "databaseAddress",
						Aliases: []string{"da"},
						Usage:   "-da 127.0.0.1:3306",
						Destination: &databaseAddress,
					},
					&cli.StringFlag{
						Name:    "databaseUser",
						Aliases: []string{"du"},
						Usage:   "-du name",
						Destination: &databaseUser,
					},
					&cli.StringFlag{
						Name:    "databasePassword",
						Aliases: []string{"dp"},
						Usage:   "-dp passwd",
						Destination: &databasePassword,
					},
					&cli.StringFlag{
						Name:    "proxyAddress",
						Aliases: []string{"pa"},
						Usage:   "-pa 127.0.0.1:8080",
						Destination: &proxyAddress,
					},
					&cli.StringFlag{
						Name:    "proxyUser",
						Aliases: []string{"pu"},
						Usage:   "-pu name",
						Destination: &proxyUser,
					},
					&cli.StringFlag{
						Name:    "proxyPassword",
						Aliases: []string{"pp"},
						Usage:   "-pp passwd",
						Destination: &proxyPassword,
					},
					&cli.StringFlag{
						Name:    "fileInput",
						Aliases: []string{"f"},
						Usage:   "-f filename(like: -f test.txt)",
						Destination: &filename,
					},
					&cli.BoolFlag{
						Name:    "WindowsAuth",
						Aliases: []string{"WA"},
						Usage:   "-WA (Only for mssql, if choose this, it will connect mssql using windows authentication)",
						Value: false,
						Destination: &windowsAuth,
					},

					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Value:   3,
						Usage:   "-n (The number for extract information from databases)",
						Destination: &num,
					},
					&cli.IntFlag{
						Name:    "thread",
						Aliases: []string{"t"},
						Value:   5,
						Usage:   "-t 1 (Only For SearchSensitiveData function)",
						Destination: &thread,
					},
					&cli.StringFlag{
						Name:    "pattern",
						Aliases: []string{"p"},
						Usage:   "-p pattern(Only For SearchSensitiveData function,like searching for username: -p ^[\\x{4e00}-\\x{9fa5}]{2,4}$ )",
						Destination: &pattern,
					},
				},
				Action:  func(c *cli.Context) error {
					defer func(){
						r := recover()
						if r != nil {
							fmt.Println("Please input the valid regular expression!")
							log.Fatal("PANIC :", r)
						}
					}()

					var start = time.Now()
					if pattern!=""{
						regexp2.MustCompile(pattern, 0)
					}
					iniInfo:=initData()
					//establish proxy connection
					var connection  net.Conn
					if proxyAddress!=""&&filename==""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
						connection= ProxyConnect(databaseAddress)
					}else if proxyAddress!=""&&filename!=""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
					}
					defer func(){
						if connection != nil {
							_ = connection.Close()
						}
					}()

					if windowsAuth{
						if databaseAddress==""{
							databaseAddress="127.0.0.1:1433"
						}
						user,_ := user.Current()
						userName:=strings.Split(user.Username, "\\")
						databaseUser=userName[len(userName)-1]
					}

					if filename==""{
						OutputID:=utils.OutputFileName(databaseAddress,databaseUser,databaseType)
						SingleSearchSensitiveData(OutputID,pattern,iniInfo)
					}else{
						inputs:=Batch(filename)
						for n,j:=range inputs{
							var connectionTmp net.Conn
							iniInfoB:=models.InitData{DatabaseType: j.Schema,DatabaseAddress: j.Address,DatabaseUser: j.User,DatabasePassword: j.Passwd,AuthSource: j.AuthSource}
							fmt.Print("No."+strconv.Itoa(n+1)+" ")
							OutputID:=utils.OutputFileName(j.Address,j.User,j.Schema)
							//establish proxy connection
							if proxyAddress!="" {
								connectionTmp = ProxyConnect(j.Address)
							}
							SingleSearchSensitiveData(OutputID,pattern,iniInfoB)
							if connectionTmp!=nil{
								err:=connectionTmp.Close()
								if err!=nil{
									log.Fatal(err)
								}
							}
						}
					}
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end)
					return nil
				},
			},

			{
				Name:    "SingleTable",
				Aliases: []string{"ST"},
				Usage:   "Command for getting data from the specified table",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "databaseType",
						Aliases: []string{"T"},
						Usage:   "-T mysql (currently supports mysql,mssql,oracle and mongo)",
						Destination: &databaseType,
					},
					&cli.StringFlag{
						Name:    "databaseAddress",
						Aliases: []string{"da"},
						Usage:   "-da 127.0.0.1:3306)",
						Destination: &databaseAddress,
					},
					&cli.StringFlag{
						Name:    "databaseUser",
						Aliases: []string{"du"},
						Usage:   "-du name",
						Destination: &databaseUser,
					},
					&cli.StringFlag{
						Name:    "databasePassword",
						Aliases: []string{"dp"},
						Usage:   "-dp passwd",
						Destination: &databasePassword,
					},
					&cli.StringFlag{
						Name:    "proxyAddress",
						Aliases: []string{"pa"},
						Usage:   "-pa 127.0.0.1:8080",
						Destination: &proxyAddress,
					},
					&cli.StringFlag{
						Name:    "proxyUser",
						Aliases: []string{"pu"},
						Usage:   "-pu name",
						Destination: &proxyUser,
					},
					&cli.StringFlag{
						Name:    "proxyPassword",
						Aliases: []string{"pp"},
						Usage:   "-pp passwd",
						Destination: &proxyPassword,
					},
					&cli.StringFlag{
						Name:    "databaseTable",
						Aliases: []string{"dt"},
						Usage:   "-dt database.table",
						Destination: &singleTable,
					},
					&cli.BoolFlag{
						Name:    "WindowsAuth",
						Aliases: []string{"WA"},
						Usage:   "-WA (Only for mssql, if choose this, it will connect mssql using windows authentication)",
						Value: false,
						Destination: &windowsAuth,
					},

					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Value:   3,
						Usage:   "-n (The number for extract information from databases)",
						Destination: &num,
					},
				},
				Action:  func(c *cli.Context) error {
					var start = time.Now()
					var connection  net.Conn
					var tableList []string
					iniInfo:=initData()
					//proxy code
					if proxyAddress!=""{
						proxyConnection=ProxyConfig(proxyAddress,proxyUser,proxyPassword)
						connection= ProxyConnect(databaseAddress)
					}
					defer func(){
						if connection != nil {
							_ = connection.Close()
						}
					}()

					if singleTable==""{
						log.Fatalf("Please input the specified table, like: -dt databaseName.tableName/collectionName")
					}else{
						if windowsAuth{
							if databaseAddress==""{
								databaseAddress="127.0.0.1:1433"
							}
							user,_ := user.Current()
							userName:=strings.Split(user.Username, "\\")
							databaseUser=userName[len(userName)-1]
						}

						tableList=append(tableList,singleTable)
						outputID:=utils.OutputFileName(databaseAddress,databaseUser,databaseType)
						if iniInfo.DatabaseType=="mongo"{
							client:=MongodbInit(iniInfo)
							SampledataMongo(client,tableList,num,outputID)
							defer client.Disconnect(context.Background())
						}else{
							db:=DBinit(iniInfo)
							Sampledata(db,tableList,num,outputID,iniInfo.DatabaseType)
						}
					}
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end)
					return nil
				},
			},

		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

/*
  overview function entrance
  @Param  outputID (the output file name)
  @Param info (the information user inputs)
*/
func SingleOverview(outputID utils.InfoStruct,info models.InitData){
	if info.DatabaseType=="mongo"{
		client:=MongodbInit(info)
		OverviewMongo(client,outputID)
		defer client.Disconnect(context.Background())
	}else{
		db:=DBinit(info)
		Overview(db,outputID,info.DatabaseType)
	}
}

/*
  sample data function entrance
  @Param  outputID (the output file name)
  @Param info (the information user inputs)
*/
func SingleSampleData(outputID utils.InfoStruct,info models.InitData){
	var tableList []string
	if info.DatabaseType=="mongo"{
		client:=MongodbInit(info)
		tableList=CountAllCollections(client)
		SampledataMongo(client,tableList,num,outputID)
		defer client.Disconnect(context.Background())
	}else{
		db:=DBinit(info)
		if info.DatabaseType=="mysql"{
			tableList=CountAllTables(db)
		}else if info.DatabaseType=="mssql"{
			tableList=CountAllTablesMs(db)
		}else if info.DatabaseType=="oracle"{
			tableList=CountAllTablesOra(db)
		}
		Sampledata(db,tableList,num,outputID,info.DatabaseType)
	}
}

/*
  search sensitive data function entrance
  @Param  outputID (the output file name)
  @Param pattern (the regular expresstion pattern)
  @Param info (the information user inputs)
*/
func SingleSearchSensitiveData(outputID utils.InfoStruct,pattern string,info models.InitData){

	var tableList []string
	if info.DatabaseType=="mongo"{
		client:=MongodbInit(info)
		tableList=CountAllCollections(client)
		LookforSensitiveData(nil,client, tableList, num, thread, outputID, pattern, info.DatabaseType)
		defer client.Disconnect(context.Background())
	}else {
		db:=DBinit(info)
		if info.DatabaseType == "mysql" {
			tableList = CountAllTables(db)
		} else if info.DatabaseType == "mssql" {
			tableList = CountAllTablesMs(db)
		} else if info.DatabaseType == "oracle" {
			tableList = CountAllTablesOra(db)
		}
		LookforSensitiveData(db,nil, tableList, num, thread, outputID, pattern, info.DatabaseType)
	}

}

/*
  structuralize the database information
  @Return model.InitData struct (database information struct)
*/
func initData() models.InitData{
	var authSource string
	if databaseType=="mongo"{
		if strings.Contains(databaseAddress,"?"){
			value:=strings.Split(databaseAddress,"?")
			if len(value)==2{
				databaseAddress=strings.Trim(value[0],"/")
				authSource=value[1]
			}else{
				log.Fatal("Please provide database name after database address, like: 127.0.0.1:27017?databaseName")
			}
		}
		return models.InitData{DatabaseType: databaseType,DatabaseAddress: databaseAddress,DatabaseUser: databaseUser,DatabasePassword: databasePassword,ProxyAddress: proxyAddress,ProxyUser: proxyUser,ProxyPassword: proxyPassword,AuthSource: authSource}
	}
	return models.InitData{DatabaseType: databaseType,DatabaseAddress: databaseAddress,DatabaseUser: databaseUser,DatabasePassword: databasePassword,ProxyAddress: proxyAddress,ProxyUser: proxyUser,ProxyPassword: proxyPassword,WindowsAuth: windowsAuth}
}

/*
  unified database initialization, and return database handle
  @Param  info (the information user inputs)
  @Return sql.DB (database handle)
*/
func DBinit(info models.InitData) (*sql.DB){
	var db *sql.DB
	switch info.DatabaseType {
	      case "mysql":
	      	db=MysqlDBinit(info)
          case "mssql":
          	db=MssqlDBinit(info)
	      case"oracle":
		    db=OracleDBinit(info)
	      default:
			  log.Fatal("Currently only supports mysql,mssql and oracle!")
	}
	return  db
}