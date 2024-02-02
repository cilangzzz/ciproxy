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
)

// Proxy Config
const (
	ProxyVersion = "v0.0.0"
	ProxyMode    = "Debug"

	ProxyOrganization = "www.cilang.buzz"
)

// connectConfig Connect Config Constant
const (
	DefaultConnectProtocol = "Tcp"
	DefaultOutTime         = 10 * time.Second
)
