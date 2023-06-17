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

type Handle func(conn net.Conn, req []byte) []byte

type (
	MiddleManage struct {
		HandleChain []Handle
	}
)

func (m MiddleManage) Use() {

}
func (m MiddleManage) Add() {

}
