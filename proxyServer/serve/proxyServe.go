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
	"fmt"
	"github.com/opencvlzg/ciproxy/constants/proxyConfig"
	"github.com/opencvlzg/ciproxy/constants/proxyMethod"
	"github.com/opencvlzg/ciproxy/proxyServer/middleHandle"
	"github.com/opencvlzg/ciproxy/proxyServer/serverHandle"
	"github.com/opencvlzg/ciproxy/util"
	"log"
)

type ProxyServe struct {
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Method   string `json:"method,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Host     string `json:"host,omitempty"`
	LogPath  string `json:"logPath,omitempty"`
}

// init 初始化
func (p ProxyServe) init() {
	util.LogInit(p.LogPath)
	p.Host = p.Ip + ":" + p.Port
	p.printInfo()
}

// Start 服务启动主入口
func (p *ProxyServe) Start() {
	p.init()
	switch p.Method {
	case proxyMethod.HttpProxy:
		p.HttpProxyListen()
	case proxyMethod.HttpsProxy:
		p.HttpsProxyListen()
	case proxyMethod.HttpsSniffProxy:
		p.HttpsSniffProxyListen()
	case proxyMethod.TcpTunnelProxy:
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

// printInfo 打印信息
func (p ProxyServe) printInfo() {
	// 3d logo
	println()
	println("/\\  _`\\    __/\\  _`\\                                \n\\ \\ \\/\\_\\ /\\_\\ \\ \\L\\ \\_ __   ___   __  _  __  __    \n \\ \\ \\/_/_\\/\\ \\ \\ ,__/\\`'__\\/ __`\\/\\ \\/'\\/\\ \\/\\ \\   \n  \\ \\ \\L\\ \\\\ \\ \\ \\ \\/\\ \\ \\//\\ \\L\\ \\/>  </\\ \\ \\_\\ \\  \n   \\ \\____/ \\ \\_\\ \\_\\ \\ \\_\\\\ \\____//\\_/\\_\\\\/`____ \\ \n    \\/___/   \\/_/\\/_/  \\/_/ \\/___/ \\//\\/_/ `/___/> \\\n                                              /\\___/\n                                              \\/__/ ")
	fmt.Printf("CiProxy Version %s, Mode %s\n", proxyConfig.ProxyVersion, proxyConfig.ProxyMode)
	println("log path setting to" + p.LogPath)
	fmt.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n", p.Ip, p.Port, p.Method, p.Protocol)
	log.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n\n", p.Ip, p.Port, p.Method, p.Protocol)

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
