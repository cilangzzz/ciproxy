/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/24
  @desc: //TODO
**/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/opencvlzg/ciproxy"
	"io"
	"net/http"
)

var buffer = bytes.NewBuffer([]byte{})
var writer = io.MultiWriter(buffer)
var buf = bufio.NewReader(buffer)
var client = &http.Client{}

func logger() {

	for {
		request, err := http.ReadRequest(buf)

		if err != nil {
			//return
		} else {
			println("[CiProxy]logger middleware", request.Method, request.Host, buf.Size())
			// 抢先请求
			//res, err := http.DefaultClient.Do(request)
			//if err != nil {
			//	println("client err")
			//	log.Println(err)
			//	//return
			//} else {
			//	println("client do")
			//	log.Println(res)
			//}

		}
	}
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "Server Ip Address")
	port := flag.String("port", "6677", "Server Port")
	method := flag.String("method", ciproxy.HttpsSniffDetailProxy, "Server METHOD NORMAL,TUNNEL, SNIFF")
	//protocol := flag.String("protocol", "TCP", "Connect Protocol")
	//logPath := flag.String("log", "log/proxy.log", "log file path")
	flag.Parse()
	//file, err := os.Create("data.txt")
	//if err != nil {
	//	return
	//}
	//defer file.Close()
	//ciproxy.DefaultWriter = bufio.NewWriter(file)
	//ciproxy.DefaultWriter = writer

	proxyServe := ciproxy.Default()
	proxyServe.Ip = *ip
	proxyServe.Port = *port
	proxyServe.Method = *method
	//go logger()
	proxyServe.Start()
}
