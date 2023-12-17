/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/12
  @desc: //TODO
**/

package localProxy

import (
	"log"
	"net"
	"trafficForward/internal/proxyClient/serverProxy"
	"trafficForward/internal/proxyClient/trafficForward"
)

type (
	LocalProxy struct {
		Ip          string
		Port        string
		Method      string
		Client      net.Conn
		ProxyClient serverProxy.ServerProxy
	}
)

func (l *LocalProxy) Start() {
	host := l.Ip + ":" + l.Port
	ln, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := ln.Accept()
		l.Client = client
		if err != nil {
			log.Println(err)
		}
		host := l.ProxyClient.Ip + ":" + l.ProxyClient.Port
		go trafficForward.HandleServerConnect(client, host, l.ProxyClient.TLSConfig)

	}

}
