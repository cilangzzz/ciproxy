package ciproxy

import "time"

// ProxyMethod constant proxyMethod
const (
	HttpProxy       = "HttpProxy"
	HttpsProxy      = "HttpsProxy"
	HttpsSniffProxy = "HttpsSniffProxy"
	CustomProxy     = "CustomProxy"
	TcpNormalProxy  = "TcpNormal"
	TcpTunnelProxy  = "TcpTunnel"
	PortProxy       = "PortProxy"
	DefaultProxy    = "All"
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
