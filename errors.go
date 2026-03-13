package ciproxy

import "errors"

// 标准错误定义
var (
	ErrInvalidRequest       = errors.New("invalid request")
	ErrConnectionFailed     = errors.New("remote connection failed")
	ErrTLSHandshakeFailed   = errors.New("TLS handshake failed")
	ErrCertificateLoad      = errors.New("certificate load failed")
	ErrProtocolNotSupported = errors.New("protocol not supported")
	ErrEntryNotFound        = errors.New("traffic entry not found")
)

// ProxyError 表示代理相关错误
type ProxyError struct {
	Op   string // 失败的操作
	Host string // 目标主机
	Err  error  // 底层错误
}

func (e *ProxyError) Error() string {
	if e.Host != "" {
		return e.Op + ": " + e.Host + ": " + e.Err.Error()
	}
	return e.Op + ": " + e.Err.Error()
}

func (e *ProxyError) Unwrap() error {
	return e.Err
}

// NewProxyError 创建新的代理错误
func NewProxyError(op, host string, err error) *ProxyError {
	return &ProxyError{Op: op, Host: host, Err: err}
}
