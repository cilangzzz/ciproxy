/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package serve

import (
	"crypto/tls"
	"log"
	"net"
	"trafficForward/server/trafficForward"
	"trafficForward/server/util"
)

type ProxyServe struct {
	Ip            string `json:"ip,omitempty"`
	Port          string `json:"port,omitempty"`
	Method        string `json:"method,omitempty"`
	ListenAddress string `json:"listen_address,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

func (p *ProxyServe) StartInTls() {
	tlsConfig := util.TLSUtil{Organization: "cilang"}
	cert, err := tlsConfig.GenCertificate()
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", p.ListenAddress, config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go trafficForward.HandleClientConnect(client)
	}
}
func (p *ProxyServe) Start() {
	ln, err := net.Listen("tcp", p.ListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		go trafficForward.HandleClientConnect(client)
	}
}
