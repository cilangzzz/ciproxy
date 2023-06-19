/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/17
  @desc: //TODO
**/

package middle

import "net"

var MdManage MiddleManage

type Handle func(conn net.Conn) []byte

type (
	MiddleManage struct {
		HandleChain []Handle
	}
)

func (m MiddleManage) Use(handle Handle) {
	m.HandleChain = append(m.HandleChain, handle)
}

//func (m MiddleManage) Add() {
//
//}
