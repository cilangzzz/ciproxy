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
	"fmt"
	"log"
	"trafficForward/internal/proxyClient/localProxy"
	"trafficForward/internal/proxyClient/serverProxy"
	util2 "trafficForward/internal/proxyClient/util"
)

func main() {
	ip := flag.String("ip", "66.151.208.210", "Proxy Server Ip Address")
	port := flag.String("port", "80", "Proxy Server Port")
	method := flag.String("method", "TCP", "Proxy Server method")
	localIp := flag.String("l_ip", "127.0.0.1", "Proxy Server Ip Address")
	localPort := flag.String("l_port", "8000", "Local Proxy Port")
	localMethod := flag.String("l_method", "TUNNEL", "Local Proxy Method")
	configEnable := flag.String("config", "NORMAL", "")
	flag.Parse()
	fmt.Printf("proxy server %s:%s method %s\nproxy local %s:%s method%s\n", *ip, *port, *method, *localIp, *localPort, *localMethod)
	err := util2.SetProxy(*localIp + ":" + *localPort)
	if err != nil {
		log.Println("Auto set proxy failed, please set proxy manually")
	}
	if *configEnable == "NORMAL" {
		log.Println("using default config file")
	}
	tlsUtil := util2.TLSUtil{Organization: "client"}
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
		ProxyClient: serverProxy.ServerProxy{
			Ip:        *ip,
			Port:      *port,
			Method:    *method,
			TLSConfig: conf,
		},
	}
	localServe.Start()

}
