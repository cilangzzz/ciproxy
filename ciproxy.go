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
)

// DefaultWriter reference gin
var DefaultWriter io.Writer = os.Stdout

//// NewProxyServe 返回服务实例
//func NewProxyServe() *serve.ProxyServe {
//	return &serve.ProxyServe{}
//}
