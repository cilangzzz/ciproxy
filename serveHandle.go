/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: //TODO
**/

package ciproxy

import (
	"log"
	"net"
)

func errLog(msg string, err error) {
	log.Println("serverHandle:" + msg + " err:" + err.Error())
}

//// baseProxyServer 基础代理监听
//func baseProxyServer(host string, proxyHandle func(c net.Conn)) {
//	ln, err := net.Listen("tcp", host)
//	if err != nil {
//		errLog("listen serve launch failed ", err)
//		return
//	}
//	for {
//		c, err := ln.Accept()
//		if err != nil {
//			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
//			return
//		}
//		go proxyHandle(c)
//	}
//}

//// HttpProxyServer normal http server 普通的http代理服务
//func HttpProxyServer(host string) {
//	ln, err := net.Listen("tcp", host)
//	if err != nil {
//		errLog("listen serve launch failed ", err)
//		return
//	}
//	for {
//		c, err := ln.Accept()
//		if err != nil {
//			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
//			return
//		}
//		go HttpProxyHandle(c)
//	}
//
//}

//// HttpsProxyServer https server Https代理服务
//func HttpsProxyServer(host string) {
//	ln, err := net.Listen("tcp", host)
//	if err != nil {
//		errLog("listen serve launch failed ", err)
//		return
//	}
//	for {
//		c, err := ln.Accept()
//		if err != nil {
//			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
//			return
//		}
//		go HttpsProxyHandle(c)
//	}
//}
//
//// HttpsSniffProxyServer https 中间人欺骗
//func HttpsSniffProxyServer(host string) {
//	ln, err := net.Listen("tcp", host)
//	if err != nil {
//		errLog("listen serve launch failed ", err)
//		return
//	}
//	for {
//		c, err := ln.Accept()
//		if err != nil {
//			errLog("connect client failed"+c.RemoteAddr().String()+"err", err)
//			return
//		}
//		go HttpsSniffProxyHandle(c)
//	}
//}
//
//// TunnelProxyServer 加密代理监听
//func TunnelProxyServer(host string) {
//	tlsCnf := &tls.Config{InsecureSkipVerify: false}
//	ln, err := tls.Listen("tcp", host, tlsCnf)
//	if err != nil {
//		errLog("listen serve launch failed ", err)
//		return
//	}
//	for {
//		tlsC, err := ln.Accept()
//		if err != nil {
//			errLog("connect client failed"+tlsC.RemoteAddr().String()+"err", err)
//			return
//		}
//		go TunnelProxyHandle(tlsC)
//	}
//}

// CustomProxyServer 自定义监听
func CustomProxyServer(p *ProxyServe) {
	ln, err := net.Listen("tcp", p.Host)
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
		ctx := p.contextPool.Get().(*Context)
		for i := 0; i < ctx.index; i++ {
			handler := ctx.handlers[i]
			handle(handler)
		}
	}

}

// handle connHandle
func handle(proxyHandle ProxyHandle) {

}
