/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/12
  @desc: //TODO
**/

package proxyClient

import (
	"crypto/tls"
	"net"
	"trafficForward/client/trafficForward"
)

var proxyClient *ProxyClient

type (
	ProxyClient struct {
		Ip        string
		Port      string
		Method    string
		Client    net.Conn
		TLSConfig *tls.Config
	}
)

func (p *ProxyClient) Connect() {
	host := proxyClient.Ip + ":" + proxyClient.Port
	trafficForward.HandleServerConnect(p.Client, host)
}
