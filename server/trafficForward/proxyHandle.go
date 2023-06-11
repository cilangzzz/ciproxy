/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package trafficForward

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"trafficForward/server/util"
)

func HandleClientConnect(client net.Conn) {
	buf := make([]byte, 1024)
	_, err := client.Read(buf)
	var method, url string
	fmt.Printf("%s", buf)
	_, err = fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &url)
	if err != nil {
		return
	}
	host := util.GetHttpsHostRegex(url)
	println(host)
	println(method)
	switch method {
	case "CONNECT":
		handleHttps(host, client)
	default:
		handleHttp(host, client)
	}

}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHttps(host string, client net.Conn) {

	target, err := net.DialTimeout("tcp", host, 60*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	client.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	go transfer(target, client)
	go transfer(client, target)
}

func handleHttp(host string, client net.Conn) {
	target, err := net.DialTimeout("tcp", host, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	go transfer(target, client)
	go transfer(client, target)
}