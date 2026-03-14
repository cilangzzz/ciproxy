/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/21
  @desc: 服务监听处理
**/

package ciproxy

import (
	"log"
	"net"
)

// ServeProxy 启动监听（向后兼容）
// Deprecated: 使用 ProxyServe.Start() 代替
func ServeProxy(p *ProxyServe) {
	// 向后兼容：如果使用旧的 API，初始化必要字段
	if p.config == nil {
		p.config = &DefaultConfig
	}
	if p.logger == nil {
		p.logger = GetLogger()
	}

	// 迁移旧字段
	p.migrateLegacyFields()

	// 初始化处理器链（如果没有初始化）
	if len(p.handlersChain) == 0 {
		p.initHandlers()
	}

	// 初始化上下文池
	p.contextPool.New = func() interface{} {
		return p.newContext()
	}

	// 初始化日志
	logInit(p.config.LogPath)

	// 启动监听
	addr := p.config.IP + ":" + p.config.Port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("listen serve launch failed ", err)
		return
	}
	p.listener = ln
	p.running = true
	p.stats.StartTime = time.Now()

	p.printBanner()

	for {
		c, err := ln.Accept()
		if err != nil {
			if p.shuttingDown {
				return
			}
			log.Println("connect client failed "+c.RemoteAddr().String()+" err", err)
			continue
		}

		// 检查连接数限制
		if p.maxConnections > 0 && p.stats.ActiveConnections >= p.maxConnections {
			log.Println("max connections reached, rejecting:", c.RemoteAddr())
			c.Close()
			continue
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
	// 预留扩展点
}