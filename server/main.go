/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/5/21
  @desc: //TODO
**/

package main

import (
	"flag"
	"fmt"
	"trafficForward/server/serve"
)

func main() {
	ip := flag.String("IP", "127.0.0.1", "Server Ip Address")
	port := flag.String("PORT", "8000", "Server Port")
	method := flag.String("METHOD", "NORMAL", "Server METHOD NORMAL,TUNNEL")
	protocol := flag.String("PROTOCOL", "TCP", "Connect Protocol")
	proxyServe := serve.ProxyServe{
		Ip:            *ip,
		Port:          *port,
		Method:        *method,
		ListenAddress: fmt.Sprintf("%s:%s", *ip, *port),
		Protocol:      *protocol,
	}
	proxyServe.Start()
}
