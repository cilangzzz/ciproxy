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
	"log"
	"net"
	"trafficForward/server/trafficForward"
)

type ProxyServe struct {
	Ip              string `json:"ip,omitempty"`
	Port            string `json:"port,omitempty"`
	IsEncryptTunnel int    `json:"isEncryptTunnel,omitempty"`
}

func (p *ProxyServe) Start() {
	ln, err := net.Listen("tcp", ":8000")
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
