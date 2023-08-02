package core

import (
	"log"
	"net"
	"golang.org/x/net/proxy"
)

/*
  fill in the proxy authentication struct, user, passwd
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
  connect the destination address
  @Param  address (proxy address)
  @Return net.Conn (the connection with proxy server)
*/
func ProxyConnect(address string) (net.Conn) {
	conn,err:=proxyConnection.Dial("tcp",address)
	if err!=nil{
		log.Fatal(err.Error())
	}
	return conn
}