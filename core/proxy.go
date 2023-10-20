package core

import (
	"log"
	"net"
	"strings"
	"golang.org/x/net/proxy"
)

/*
  Fill in the proxy authentication struct, user, passwd
  @Param  pa (proxy address)
  @Param  pu (proxy name)
  @Param  pp (proxy password)
  @Return proxy.Dialer (the dialer connection with proxy server)
*/
func ProxyConfig(pa,pu,pp string) (proxy.Dialer) {
	auth:= proxy.Auth {pu,pp}
	dialer,err:=proxy.SOCKS5("tcp",pa,&auth,proxy.Direct)
	if err!=nil{
		log.Fatal(err.Error())
	}
	return dialer
}

/*
  Connect the destination address
  @Param  address (proxy address)
  @Return net.Conn (the connection with proxy server)
*/
func ProxyConnect(address string) (net.Conn) {
	parts := strings.Split(address, "/")
	conn,err:=proxyConnection.Dial("tcp",parts[0])
	if err!=nil{
		log.Fatal(err.Error())
	}
	return conn
}

/*
  Connect the destination address, if something goes wrong, just return the error and don't exit
  @Param  address (proxy address)
  @Return net.Conn (the connection with proxy server)
  @Return error
*/
func ProxyConnectNoExit(address string) (net.Conn,error) {
	parts := strings.Split(address, "/")
	conn,err:=proxyConnection.Dial("tcp",parts[0])
	if err!=nil{
		return nil,err
	}
	return conn,nil
}