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

// ServeProxy 启动监听
func ServeProxy(p *ProxyServe) {
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
		// 获取上下文
		ctx := p.contextPool.Get().(*Context)
		// 设置上下文客户端
		ctx.ClientConn = c
		// 响应链处理
		ctx.Next()
		// 重置上下文
		ctx.reset()
		// 放回池
		p.contextPool.Put(ctx)
	}

}

// handle connHandle
func handle(proxyHandle ProxyHandle) {

}
