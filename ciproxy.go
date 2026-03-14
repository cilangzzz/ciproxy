/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2024/1/3
  @desc: //TODO
**/

// Package ciproxy a proxy frame implement by tcp,udp
package ciproxy

import (
	"io"
	"os"

	// 重新导出核心类型
	"github.com/opencvlzg/ciproxy/pkg/context"
	"github.com/opencvlzg/ciproxy/pkg/middleware"
)

// ========== 类型别名，保持向后兼容 ==========

// Context 请求上下文
type Context = context.Context

// ProxyHandle 代理处理函数
type ProxyHandle = context.ProxyHandle

// ProxyHandlersChain 处理器链
type ProxyHandlersChain = context.ProxyHandlersChain

// Middleware 中间件接口
type Middleware = middleware.Middleware

// MiddlewareFunc 中间件函数类型
type MiddlewareFunc = middleware.MiddlewareFunc

// ========== 包级导出 ==========

// DefaultWriter reference gin
var DefaultWriter io.Writer = os.Stdout

// NewProxyServe 返回服务实例
func NewProxyServe() *ProxyServe {
	return &ProxyServe{}
}

// Default 返回默认服务实例
func Default() *ProxyServe {
	return &ProxyServe{
		Ip:                 DefaultIp,
		Port:               DefaultPort,
		Method:             DefaultProxy,
		Protocol:           DefaultConnectProtocol,
		LogPath:            "",
		Host:               "",
		ProxyHandlersChain: nil,
	}
}
