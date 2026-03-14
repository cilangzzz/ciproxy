package ciproxy

import "time"

// ProxyMethod constant proxyMethod
const (
	HttpProxy             = "HttpProxy"
	HttpsProxy            = "HttpsProxy"
	HttpsSniffProxy       = "HttpsSniffProxy"
	HttpsSniffDetailProxy = "HttpsSniffDetailProxy"
	HttpInterceptProxy    = "HttpInterceptProxy" // 新增：完整的 HTTPS MITM 拦截代理
	WebsocketProxy        = "WebsocketProxy"
	TcpNormalProxy        = "TcpNormal"
	TcpTunnelProxy        = "TcpTunnel"
	PortProxy             = "PortProxy"
	DefaultProxy          = "All"
)

// Proxy Config
const (
	ProxyVersion = "v0.0.0"
	ProxyMode    = "Debug"

	// DefaultIp DefaultPort defaultServerConfig
	DefaultIp   = "127.0.0.1"
	DefaultPort = ""

	ProxyOrganization = "www.cilang.buzz"
)

// connectConfig Connect Config Constant
const (
	DefaultConnectProtocol = "Tcp"
	DefaultOutTime         = 10 * time.Second
)
