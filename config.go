package ciproxy

import (
	"encoding/json"
	"os"
	"time"
)

// TLSConfig TLS证书配置
type TLSConfig struct {
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
	CertData []byte `json:"certData,omitempty"`
	KeyData  []byte `json:"keyData,omitempty"`
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	Connect time.Duration `json:"connect"`
	Read    time.Duration `json:"read"`
	Write   time.Duration `json:"write"`
}

// FeatureFlags 功能开关
type FeatureFlags struct {
	EnableTrafficCapture bool `json:"enableTrafficCapture"`
	EnableTrafficReplay  bool `json:"enableTrafficReplay"`
	EnableTunnelCrypt    bool `json:"enableTunnelCrypt"`
}

// ProxyConfig 代理服务配置
type ProxyConfig struct {
	// 网络设置
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`

	// 代理设置
	Method  string `json:"method"`
	LogPath string `json:"logPath"`

	// TLS设置
	TLS TLSConfig `json:"tls"`

	// 超时设置
	Timeout TimeoutConfig `json:"timeout"`

	// 功能开关
	Features FeatureFlags `json:"features"`
}

// DefaultConfig 默认配置
var DefaultConfig = ProxyConfig{
	IP:       DefaultIp,
	Port:     "8080",
	Protocol: DefaultConnectProtocol,
	Method:   HttpProxy,
	LogPath:  "",
	TLS: TLSConfig{
		CertPath: "./cert/root.crt",
		KeyPath:  "./cert/private.pem",
	},
	Timeout: TimeoutConfig{
		Connect: DefaultOutTime,
		Read:    30 * time.Second,
		Write:   30 * time.Second,
	},
	Features: FeatureFlags{
		EnableTrafficCapture: false,
		EnableTrafficReplay:  false,
		EnableTunnelCrypt:    false,
	},
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*ProxyConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save 保存配置到文件
func (c *ProxyConfig) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Validate 验证配置有效性
func (c *ProxyConfig) Validate() error {
	if c.IP == "" {
		c.IP = DefaultIp
	}
	if c.Port == "" {
		c.Port = "8080"
	}
	if c.Method == "" {
		c.Method = HttpProxy
	}
	if c.Timeout.Connect == 0 {
		c.Timeout.Connect = DefaultOutTime
	}
	return nil
}
