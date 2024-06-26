/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2024/1/10
  @desc: // reference from gin context design
**/

package ciproxy

import (
	"net"
	"net/http"
	"sync"
)

// ProxyHandle define the proxyHandle to handle proxyRequest and used by middleware
type ProxyHandle func(ctx *Context)

// ProxyHandlersChain define proxyHandle slice
type ProxyHandlersChain []ProxyHandle

// Context reference from go-proxy and gin
type Context struct {
	// ConnStatus conn status connecting closed nil
	ConnStatus string
	// ClientConn client net conn
	ClientConn net.Conn
	// TlsClientConn client tls net conn
	TlsClientConn net.Conn
	Request       *http.Request

	index    int
	handlers ProxyHandlersChain

	// protect middleware context
	mu sync.RWMutex

	// ServerConn server net conn
	ServerConn net.Conn
	// TslServerConn server tls net conn
	TslServerConn net.Conn
	Response      http.Response
}

// SetClientConn set client conn
func (c *Context) SetClientConn(ClientConn net.Conn) {
	c.ClientConn = ClientConn
}

// SetRequest set the request to the context.Request
func (c *Context) SetRequest() {

}

// SetServerConn set the server conn
func (c *Context) SetServerConn(ServerConn net.Conn) {
	c.ServerConn = ServerConn
}

// SetResponse set the response to the context.Response
func (c *Context) SetResponse() {

}

// IsAbort return abort label
func (c *Context) IsAbort() bool {
	return c.index == -1
}

// Abort set the abort index
func (c *Context) Abort() {
	c.index = -1
}

// Next set to next handle
func (c *Context) Next() {
	//c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// reset reset the context
func (c *Context) reset() {
	c.index = 0
	c.ConnStatus = "closed"
	c.ClientConn = nil
	c.Request = nil
	//c.handlers = nil
	c.ServerConn = nil
	c.Response = http.Response{}

}

//func (c *Context)GetTransport()    *http.Transport{
//	return &http.Transport{DialTLS: c.TslServerConn.RemoteAddr().String()}
//}
