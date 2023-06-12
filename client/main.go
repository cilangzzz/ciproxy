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
	"trafficForward/client/localProxy"
	"trafficForward/client/proxyClient"
)

func main() {
	ip := flag.String("IP", "66.151.208.210", "Proxy Server Ip Address")
	port := flag.String("PORT", "8000", "Proxy Server Port")
	method := flag.String("METHOD", "TCP", "Proxy Server method")
	localIp := flag.String("i", "127.0.0.1", "Proxy Server Ip Address")
	localPort := flag.String("p", "8000", "Local Proxy Port")
	localMethod := flag.String("m", "8000", "Local Proxy Method")

	_ = proxyClient.ProxyClient{
		Ip:        *ip,
		Port:      *port,
		Method:    *method,
		TLSConfig: nil,
	}

	_ = localProxy.LocalProxy{
		Ip:     *localIp,
		Port:   *localPort,
		Method: *localMethod,
	}

}
