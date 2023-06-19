/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package trafficHandle

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"trafficForward/server/middle"
	"trafficForward/server/util"
)

func HandleClientConnect(client net.Conn) {

	buf := make([]byte, 1024)
	_, err := client.Read(buf)

	host, err := util.ParseUrl(buf)
	log.Println(host)
	if err != nil {
		return
	}
	var method, url string
	_, err = fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &url)
	if err != nil {
		return
	}

	if false {
		log.Println(buf)
		return
	}

	switch method {
	case "CONNECT":
		handleHttps(host, client)
	default:
		handleHttps(host, client)
	}

}

func handleHttps(host string, client net.Conn) {

	target, err := net.DialTimeout("tcp", host, 30*time.Second)
	if err != nil {
		println(host)
		log.Println(err)
		return
	}
	_, err = client.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println(err)
		return
	}
	mdManage := middle.MdManage
	for i, handle := range mdManage.HandleChain {
		println(i)
		println(handle)
	}
	go transfer(target, client)
	go transfer(client, target)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			log.Println(err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			log.Println(err)
		}
	}(source)
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Println(err)
	}
}

func handleHttp(host string, client net.Conn) {
	target, err := net.DialTimeout("tcp", host, 60*time.Second)
	if err != nil {
		log.Println(err)
		return
	}
	go transfer(target, client)
	go transfer(client, target)
}
