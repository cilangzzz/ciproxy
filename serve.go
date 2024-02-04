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
	"log"
	"os"
	"sync"
)

type ProxyServe struct {
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Method   string `json:"method,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	LogPath  string `json:"logPath,omitempty"`
	Host     string
	// pool contextPool
	pool sync.Pool
	// ProxyHandle use custom proxy need
	ProxyHandle func(ctx *Context)
}

// init 初始化
func (p *ProxyServe) init() {
	logInit(p.LogPath)
	p.Host = p.Ip + ":" + p.Port
	p.printInfo()
}

// SetProxyHandle 设置自定义代理响应处理
func (p *ProxyServe) SetProxyHandle(proxyHandle ProxyHandle) {
	p.ProxyHandle = proxyHandle
}

// Start 服务启动主入口
func (p *ProxyServe) Start() {
	// init struct filed
	p.init()
	//c := p.pool.Get().(*Context)
	//c.reset()
	switch p.Method {
	case HttpProxy:
		p.HttpProxyListen()
	case HttpsProxy:
		p.HttpsProxyListen()
	case HttpsSniffProxy:
		p.HttpsSniffProxyListen()
	case TcpTunnelProxy:
		p.TcpTunnelTlsProxyListen()
	case CustomProxy:
		p.CustomProxyListen()
	default:
		log.Println("mainServe: No proxy method had been chose")
	}
}

// HttpProxyListen Http代理监听
func (p ProxyServe) HttpProxyListen() {
	HttpProxyServer(p.Host)
}

// HttpsProxyListen Https代理监听
func (p ProxyServe) HttpsProxyListen() {
	HttpsProxyServer(p.Host)
}

// TcpTunnelTlsProxyListen Tcp隧道加密代理监听
func (p ProxyServe) TcpTunnelTlsProxyListen() {
	TunnelProxyServer(p.Host)
}

// HttpsSniffProxyListen Https代理欺骗监听
func (p ProxyServe) HttpsSniffProxyListen() {
	HttpsSniffProxyServer(p.Host)
}

// CustomProxyListen 自定义代理监听
func (p ProxyServe) CustomProxyListen() {
	if p.ProxyHandle == nil {
		log.Println("mainServe: no custom proxyHandle")
		return
	}
	CustomProxyServer(p.Host, p.ProxyHandle)
}

func (p ProxyServe) WebsocketProxyListen() {
	// Todo - no implement
}

// printInfo 打印信息
func (p ProxyServe) printInfo() {
	// 3d logo
	println("/\\  _`\\    __/\\  _`\\                                \n\\ \\ \\/\\_\\ /\\_\\ \\ \\L\\ \\_ __   ___   __  _  __  __    \n \\ \\ \\/_/_\\/\\ \\ \\ ,__/\\`'__\\/ __`\\/\\ \\/'\\/\\ \\/\\ \\   \n  \\ \\ \\L\\ \\\\ \\ \\ \\ \\/\\ \\ \\//\\ \\L\\ \\/>  </\\ \\ \\_\\ \\  \n   \\ \\____/ \\ \\_\\ \\_\\ \\ \\_\\\\ \\____//\\_/\\_\\\\/`____ \\ \n    \\/___/   \\/_/\\/_/  \\/_/ \\/___/ \\//\\/_/ `/___/> \\\n                                              /\\___/\n                                              \\/__/ ")
	fmt.Printf("CiProxy Version %s, Mode %s\n", ProxyVersion, ProxyMode)
	println("Log path setting to " + p.LogPath)
	fmt.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n", p.Ip, p.Port, p.Method, p.Protocol)
	log.Printf("Listen on %s:%s, Proxy Method %s, Ip Protocol %s\n", p.Ip, p.Port, p.Method, p.Protocol)

}

// logInit 日记初始化
// init log
// if path equal none print to terminal
// if path not none, try to create dir
func logInit(path string) {
	//currentTime := time.Now().String()
	//log.SetPrefix("[" + currentTime + "]")
	log.SetOutput(DefaultWriter)
	log.SetPrefix("[CiProxy]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	// not setting return
	if path == "" {
		println("log path not setting, print to terminal")
		return
	}
	_, err := os.Stat(path)
	if err != nil {
		println("log path not exist auto create")
		err := os.Mkdir("./log/", os.ModePerm)
		if err != nil {
			println("create log file")
			println(err)
			panic(err)
		}
	}
	filePath := "./" + path
	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		println("log file setting err")
		panic(err)
	}
	// had set the path redirect to the path io
	log.SetOutput(logFile) // 将文件设置为log输出的文件
}
