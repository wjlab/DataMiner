package core


import (
	"log"
	"dataMiner/utils"
	"dataMiner/dblib"
	"dataMiner/models"
	"golang.org/x/net/proxy"
	"github.com/gookit/color"
)

/*
  Load the tnsnamse.ora file and get the available database address in this file
  @Param proxyConn (socks5 proxy dialer)
  @Param info (the information user inputs)
  @Return string type (available database address)
*/
func TNSAddressConnect(info *models.InitData,proxyConn *proxy.Dialer)string{
	// use tnsnames.ora to connect oracle database
	var connStrFinal string
	tnsEntries,err := utils.LoadTNS(info.TNSFile)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	for _,v:=range tnsEntries{
		if v.Service==""||len(v.Servers)==0{
			continue
		}
		for _,j:=range v.Servers {
			connStr := "oracle://" + info.DatabaseUser + ":" + info.DatabasePassword + "@" + j.Host + ":" + j.Port + "/" + v.Service
			color.Infoln("[*] Trying "+info.DatabaseUser+":"+info.DatabasePassword+"@"+j.Host+":"+j.Port+"/"+v.Service, " to connect oracle database ...")
			// use socks5 proxy to connect database
			if proxyConn != nil {
				connection, err := ProxyConnectNoExit(j.Host + ":" + j.Port)
				if err != nil {
					log.Println(err)
					continue
				}
				if connection != nil {
					if dblib.CheckConnectivity(connStr) {
						connection.Close()
						connStrFinal = connStr
						info.DatabaseAddress = j.Host + ":" + j.Port
						info.DatabaseInstance = v.Service
						return j.Host + ":" + j.Port
					}
					connection.Close()
				}
			}else{
				// connect database directly
				if dblib.CheckConnectivity(connStr) {
					connStrFinal = connStr
					info.DatabaseAddress = j.Host + ":" + j.Port
					info.DatabaseInstance = v.Service
					return j.Host + ":" + j.Port
				}
			}
		}
	}
	if connStrFinal==""{
		log.Fatalf("There is no available TNS entry in the tnsnames.ora file: "+info.TNSFile+". Please check username, password or tnsnames.ora file.")
	}
	return ""
}