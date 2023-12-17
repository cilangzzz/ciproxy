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

import (
	"log"
	"net"
	"trafficForward/internal/proxyServer/util"
)

var MdManage MiddleManage

type Handle func(client net.Conn, target net.Conn)

type (
	MiddleManage struct {
		HandleChain []Handle
	}
)

func (m MiddleManage) Add(handle Handle) {
	m.HandleChain = append(m.HandleChain, handle)
}

func (m MiddleManage) DenyOnly(host string) {
	m.Add(func(client net.Conn, target net.Conn) {
		if target.RemoteAddr().String() == host {
			err := target.Close()
			if err != nil {
				log.Println(err)
			}
			_, err = client.Write(util.GenerateHttp("<h1>Banned Host</h1>"))
			if err != nil {
				log.Println(err)
			}
		}
	})
}
func (m MiddleManage) PermitOnly(host string) {
	m.Add(func(client net.Conn, target net.Conn) {
		if target.RemoteAddr().String() != host {
			err := target.Close()
			if err != nil {
				log.Println(err)
			}
			_, err = client.Write(util.GenerateHttp("<h1>Banned Host</h1>"))
			if err != nil {
				log.Println(err)
			}
		}
	})
}
