/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/17
  @desc: //TODO
**/

// Package middleHandle 中间件处理
package middleHandle

import (
	"net"
)

var MdManage MiddleManage

type Handle func(client net.Conn, target net.Conn)

type (
	MiddleManage struct {
		HandleChain []Handle
	}
)

// MiddleHandle 中间件处理 外部使用
func MiddleHandle(c net.Conn, s net.Conn) {
	for _, handle := range MdManage.HandleChain {
		handle(c, s)
	}
}

// Add 添加中间件
func Add(handle Handle) {
	MdManage.HandleChain = append(MdManage.HandleChain, handle)
}
