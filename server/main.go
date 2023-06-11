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
	"log"
	"net"
	"trafficForward/server/trafficForward"
)

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		client, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		trafficForward.HandleClientConnect(client)
	}
}
