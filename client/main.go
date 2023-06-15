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
	"crypto/tls"
	"flag"
	"log"
	"trafficForward/client/localProxy"
	"trafficForward/client/proxyClient"
	"trafficForward/client/util"
)

func main() {
	ip := flag.String("IP", "66.151.208.210", "Proxy Server Ip Address")
	port := flag.String("PORT", "80", "Proxy Server Port")
	method := flag.String("METHOD", "TCP", "Proxy Server method")
	localIp := flag.String("lIP", "127.0.0.1", "Proxy Server Ip Address")
	localPort := flag.String("lPORT", "8000", "Local Proxy Port")
	localMethod := flag.String("lMETHOD", "TUNNEL", "Local Proxy Method")
	flag.Parse()
	tlsUtil := util.TLSUtil{Organization: "client"}
	cert, err := tlsUtil.GenCertificate()
	if err != nil {
		log.Fatal(err)

	}
	conf := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	localServe := localProxy.LocalProxy{
		Ip:     *localIp,
		Port:   *localPort,
		Method: *localMethod,
		ProxyClient: proxyClient.ProxyClient{
			Ip:        *ip,
			Port:      *port,
			Method:    *method,
			TLSConfig: conf,
		},
	}
	localServe.Start()

}
