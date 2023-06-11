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
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			sData := make([]byte, 1024)
			n, err := conn.Read(sData)
			fmt.Printf("Received data: %s\n", sData[:n])
			if err != nil {
				log.Fatal(err)
			}

		}
	}()
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			_, err = conn.Write([]byte("Hello from client!"))
			if err != nil {
				log.Fatal(err)
			}
			println("write data to server")
		}
	}()
	for {
		time.Sleep(10000 * time.Millisecond)
		println("main goroutine")
	}

}
