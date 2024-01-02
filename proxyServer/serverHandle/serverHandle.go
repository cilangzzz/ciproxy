/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: //TODO
**/

// Package serverHandle 服务监听处理
package serverHandle

import (
	"github.com/opencvlzg/ciproxy/proxyServer/proxyHandle"
	"log"
	"net"
)

func errLog(msg string, err error) {
	log.Println("serverHandle:" + msg + " err:" + err.Error())
}

// HttpProxyServer normal http server 普通的http代理服务
func HttpProxyServer(host string) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		errLog("listen serve launch failed ", err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
			return
		}
		go proxyHandle.HttpProxyHandle(c)
	}

}

// HttpsProxyServer https server Https代理服务
func HttpsProxyServer(host string) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		errLog("listen serve launch failed ", err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
			return
		}
		go proxyHandle.HttpsProxyHandle(c)
	}
}

// HttpsSniffProxyServer https 中间人欺骗
func HttpsSniffProxyServer(host string) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		errLog("listen serve launch failed ", err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
			return
		}
		go proxyHandle.HttpsSniffProxyHandle(c)
	}
}

// CustomProxyServer 自定义监听
func CustomProxyServer(host string, proxyHandle func(c net.Conn)) {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		errLog("listen serve launch failed ", err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
			return
		}
		go proxyHandle(c)
	}
}
