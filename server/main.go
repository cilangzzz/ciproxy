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
	ip := flag.String("ip", "", "Server Ip Address")
	port := flag.String("port", "80", "Server Port")
	method := flag.String("method", "TUNNEL", "Server METHOD NORMAL,TUNNEL")
	protocol := flag.String("protocol", "TCP", "Connect Protocol")
	flag.Parse()
	proxyServe := serve.ProxyServe{
		Ip:            *ip,
		Port:          *port,
		Method:        *method,
		ListenAddress: fmt.Sprintf("%s:%s", *ip, *port),
		Protocol:      *protocol,
	}
	fmt.Printf("Listen on %s:%s, Method %s, Traffic %s\n", *ip, *port, *method, *protocol)
	proxyServe.Start()
}
