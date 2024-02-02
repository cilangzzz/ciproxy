/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

package ciproxy

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
	LogPath  string `json:"logPath,omitempty"`
	Host     string
	// ProxyHandle use custom proxy need
	ProxyHandle func(ctx *Context)
}

// init 初始化
func (p *ProxyServe) init() {
	util.LogInit(p.LogPath)
	p.Host = p.Ip + ":" + p.Port
	p.printInfo()
}

// Use 使用中间件
func (p ProxyServe) Use(handle middleHandle.Handle) {
	middleHandle.Add(handle)
}

// SetProxyHandle 设置自定义代理响应处理
func (p *ProxyServe) SetProxyHandle(proxyHandle ProxyHandle) {
	p.ProxyHandle = proxyHandle
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
	case proxyMethod.CustomProxy:
		p.CustomProxyListen()
	default:
		log.Println("mainServe: No proxy method had been chose")
	}
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
	serverHandle.TunnelProxyServer(p.Host)
}

// HttpsSniffProxyListen Https代理欺骗监听
func (p ProxyServe) HttpsSniffProxyListen() {
	serverHandle.HttpsSniffProxyServer(p.Host)
}

// CustomProxyListen 自定义代理监听
func (p ProxyServe) CustomProxyListen() {
	if p.ProxyHandle == nil {
		log.Println("mainServe: no custom proxyHandle")
		return
	}
	//serverHandle.CustomProxyServer(p.Host, p.ProxyHandle)
}

func (p ProxyServe) WebsocketProxyListen() {
	// Todo - no implement
}

// printInfo 打印信息
func (p ProxyServe) printInfo() {
	// 3d logo
	println("/\\  _`\\    __/\\  _`\\                                \n\\ \\ \\/\\_\\ /\\_\\ \\ \\L\\ \\_ __   ___   __  _  __  __    \n \\ \\ \\/_/_\\/\\ \\ \\ ,__/\\`'__\\/ __`\\/\\ \\/'\\/\\ \\/\\ \\   \n  \\ \\ \\L\\ \\\\ \\ \\ \\ \\/\\ \\ \\//\\ \\L\\ \\/>  </\\ \\ \\_\\ \\  \n   \\ \\____/ \\ \\_\\ \\_\\ \\ \\_\\\\ \\____//\\_/\\_\\\\/`____ \\ \n    \\/___/   \\/_/\\/_/  \\/_/ \\/___/ \\//\\/_/ `/___/> \\\n                                              /\\___/\n                                              \\/__/ ")
	fmt.Printf("CiProxy Version %s, Mode %s\n", proxyConfig.ProxyVersion, proxyConfig.ProxyMode)
	println("Log path setting to " + p.LogPath)
	fmt.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n", p.Ip, p.Port, p.Method, p.Protocol)
	log.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n", p.Ip, p.Port, p.Method, p.Protocol)

}
