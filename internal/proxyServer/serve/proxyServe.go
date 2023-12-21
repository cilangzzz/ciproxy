/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package serve

import (
	"log"
	"trafficForward/internal/constant"
	"trafficForward/internal/proxyServer/middleHandle"
	"trafficForward/internal/proxyServer/serverHandle"
)

type ProxyServe struct {
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Method   string `json:"method,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
}

// Start 服务启动主入口
func (p *ProxyServe) Start() {
	p.Host = p.Ip + ":" + p.Port
	switch p.Method {
	case constant.HttpProxy:
		p.HttpProxyListen()
	case constant.HttpsProxy:
		p.HttpsProxyListen()
	case constant.HttpsSniffProxy:
		p.HttpsSniffProxyListen()
	case constant.TcpTunnelProxy:
		p.TcpTunnelTlsProxyListen()
	default:
		log.Println("mainServe: No proxy method had been chose")
	}
}

// Use 使用中间件
func (p ProxyServe) Use(handle middleHandle.Handle) {
	middleHandle.MdManage.Add(handle)
}

// HttpProxyListen Http代理监听
func (p ProxyServe) HttpProxyListen() {
	serverHandle.HttpProxyServer(p.Host)
}

// HttpsProxyListen Https代理监听
func (p ProxyServe) HttpsProxyListen() {
	serverHandle.HttpsProxyServer(p.Host)
}

// TcpTunnelTlsProxyListen Tcp隧道加密代理监听
func (p ProxyServe) TcpTunnelTlsProxyListen() {

}

// HttpsSniffProxyListen Https代理欺骗监听
func (p ProxyServe) HttpsSniffProxyListen() {

}

func (p ProxyServe) ListenTunnelTls() {
	//	tlsConfig := util.TLSUtil{Organization: "CiproxyOrganization"}
	//	cert, err := tlsConfig.GenCertificate()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	//	ln, err := tls.Listen("tcp", p.Ip+":"+p.Port, config)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for {
	//		client, err := ln.Accept()
	//
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		go trafficHandle.HandleClientConnect(client)
	//	}
	//}
	//func (p *ProxyServe) ListenNormalHttps() {
	//	ln, err := net.Listen("tcp", p.Ip+":"+p.Port)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	for {
	//		client, err := ln.Accept()
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		go trafficHandle.HandleClientConnect(client)
	//	}
}

func (p *ProxyServe) ListenHttpsListen() {
	//tlsConfig := util.TLSUtil{Organization: "CiproxyOrganization"}
	//cert, err := tlsConfig.GenCertificate()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	//ln, err := tls.Listen("tcp", p.Ip+":"+p.Port, config)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for {
	//	client, err := ln.Accept()
	//	go httpHandle2.HandleClientConnect(client)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}
}

func (p ProxyServe) HttpsSniffProxyServe() {
	//tlsClientConn := tls.Server()
	//defer func() {
	//	_ = tlsClientConn.Close()
	//}()
}
