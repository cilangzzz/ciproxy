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
	"github.com/opencvlzg/ciproxy"
	"log"
)

func logger(c *ciproxy.Context) {
	log.Println("[CiProxy]logger middleware", c.ClientConn.LocalAddr(), c.ConnStatus)
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "6677", "Server Port")
	method := flag.String("method", ciproxy.HttpsProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
	//protocol := flag.String("protocol", "TCP", "Connect Protocol")
	//logPath := flag.String("log", "log/proxy.log", "log file path")
	flag.Parse()
	proxyServe := ciproxy.Default()
	proxyServe.Ip = *ip
	proxyServe.Port = *port
	proxyServe.Method = *method
	proxyServe.AddMiddleware(logger)
	proxyServe.Start()
}
