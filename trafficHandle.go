/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

// Package ciproxy 流量转发
package ciproxy

import (
	"bufio"
	"fmt"
	"github.com/cilangzzz/ciproxy/internal/util"
	"io"
	"log"
	"net/http"
)

// cryptTraffic 流量加密
func cryptTraffic() {

}

// decryptTraffic 流量解密
func decryptTraffic() {

}

// TunnelTransfer 加密流量转发
func TunnelTransfer() {

}

// Transfer traffic transfer 流量Io转发
func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(source)
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Println("copy data transfer failed", err)
	}
}

// TeeTransfer traffic transfer 流量Io转发
func TeeTransfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			//log.Println("close io writer failed", err)
		}
	}(source)
	teeReader := io.TeeReader(source, DefaultWriter)
	//http.Read
	_, err := io.Copy(destination, teeReader)
	if err != nil {
		log.Println("copy data transfer failed", err)
	}
}

// TeeDoRequestTransfer traffic transfer 流量Io转发,手动处理请求
func TeeDoRequestTransfer(c *Context) {
	//defer func(destination io.WriteCloser) {
	//	err := destination.Close()
	//	if err != nil {
	//		//log.Println("close io writer failed", err)
	//	}
	//}(c.TlsClientConn)
	//defer func(source io.ReadCloser) {
	//	err := source.Close()
	//	if err != nil {
	//		//log.Println("close io writer failed", err)
	//	}
	//}(c.TslServerConn)

	//teeReader := io.TeeReader(c, DefaultWriter)
	// http.Read
	cReader := bufio.NewReader(c.TlsClientConn)
	request, err := http.ReadRequest(cReader)
	if err != nil {
		return
	}
	//println(request.Host)
	newReq, _ := util.NewRequest(request)
	// 创建 HTTP 请求, TODO 需要加入池优化性能
	// TODO 从TLS连接创建一个HTTP客户端
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		TLSClientConfig: tlsConfig,
	//		DialTLS: func(network, addr string) (net.Conn, error) {
	//			return tlsClient, nil
	//		},
	//	},
	//}
	//request.RequestURI = ""
	client := http.Client{}
	//request.RemoteAddr = c.TslServerConn.RemoteAddr().String()
	//request.URL.Scheme = "http"
	//request.URL.Host = "cn.bing.com"
	//u, err := url.Parse(request.URL.Scheme + "://" + request.Host + request.URL.String())
	//if err != nil {
	//	panic(err)
	//}
	//request.URL = u
	//c.Req = req
	newReq.URL.Scheme = "https"
	newReq.URL.Host = request.Host
	resp, err := client.Do(newReq)
	//resp, err := client.Get("https://cn.bing.com")
	if err != nil {
		fmt.Println("Error do request:", err)
		return
	}
	bytes, err := util.ResponseToBytes(resp)
	if err != nil {
		return
	}

	// 将响应内容写回 TCP 连接
	c.TlsClientConn.Write(bytes)
	//bytes := util.HttpContext("你好")
	//c.Write(bytes)

}
