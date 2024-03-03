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
	contextPool sync.Pool
	// ProxyHandle use custom proxy need
	ProxyHandlersChain ProxyHandlersChain
}

// init 初始化
func (p *ProxyServe) init() {
	logInit(p.LogPath)
	p.Host = p.Ip + ":" + p.Port
	p.printInfo()
	// sync.pool
	p.contextPool.New = func() any {
		return p.newContext()
	}
}

// newContext create new context
func (p *ProxyServe) newContext() *Context {
	return &Context{handlers: p.ProxyHandlersChain}
}

// AddHandle 设置自定义代理响应处理
// 从尾部添加响应处理
// set the custom proxyHandle if u choose the customMethod
func (p *ProxyServe) AddHandle(proxyHandle ProxyHandle) {
	p.ProxyHandlersChain = append(p.ProxyHandlersChain, proxyHandle)
}

// AddMiddleware 从头部添加中间件
// set the custom middleware
func (p *ProxyServe) AddMiddleware(proxyHandle ProxyHandle) {
	p.ProxyHandlersChain = append([]ProxyHandle{proxyHandle}, p.ProxyHandlersChain...)
}

// Start 服务启动主入口
// running the server
func (p *ProxyServe) Start() {
	// init struct filed
	p.init()

	switch p.Method {
	case HttpProxy:
		p.AddHandle(HttpProxyHandle)
	case HttpsProxy:
		p.AddHandle(HttpsProxyHandle)
	case HttpsSniffProxy:
		p.AddHandle(HttpsSniffProxyHandle)
	case TcpTunnelProxy:
		p.AddHandle(TunnelProxyHandle)
	case WebsocketProxy:
		p.AddHandle(WebsocketProxyHandle)
	case DefaultProxy:
		if len(p.ProxyHandlersChain) == 0 {
			panic("mainServe: No proxy handle implement")
		}
	default:
		panic("mainServe: No proxy method had been chose")
	}
	p.ServerHandleListen()
}

// ServerHandleListen ServerHandle 服务代理处理
func (p *ProxyServe) ServerHandleListen() {
	ServeProxy(p)
}

// printInfo 打印信息
func (p *ProxyServe) printInfo() {
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
