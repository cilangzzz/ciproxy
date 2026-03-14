/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/11
  @desc: 代理服务器核心实现，参考Gin/Echo框架设计
**/

package ciproxy

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ServerStats 服务器统计信息
type ServerStats struct {
	TotalConnections  int64     // 总连接数
	ActiveConnections int64     // 当前活跃连接数
	FailedConnections int64     // 失败连接数
	BytesReceived     int64     // 接收字节数
	BytesSent         int64     // 发送字节数
	StartTime         time.Time // 启动时间
}

// ProxyServe 代理服务器
type ProxyServe struct {
	// 配置
	config *ProxyConfig

	// 核心组件
	logger   *Logger
	listener net.Listener

	// 处理器链
	handlersChain ProxyHandlersChain
	contextPool   sync.Pool

	// 状态管理
	mu           sync.RWMutex
	running      bool
	shuttingDown bool

	// 统计信息
	stats ServerStats

	// 连接控制
	maxConnections int64
	connWg         sync.WaitGroup

	// 优雅关闭
	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc

	// 信号处理
	signalChan chan os.Signal

	// ========== 向后兼容字段 ==========
	// Deprecated: 使用 config.IP 代替
	Ip string `json:"ip,omitempty"`
	// Deprecated: 使用 config.Port 代替
	Port string `json:"port,omitempty"`
	// Deprecated: 使用 config.Method 代替
	Method string `json:"method,omitempty"`
	// Deprecated: 使用 config.Protocol 代替
	Protocol string `json:"protocol,omitempty"`
	// Deprecated: 使用 config.LogPath 代替
	LogPath string `json:"logPath,omitempty"`
	// Deprecated: 内部使用
	Host string
	// Deprecated: 内部使用
	ProxyHandlersChain ProxyHandlersChain
}

// ServerOption 服务器选项函数
type ServerOption func(*ProxyServe)

// New 创建新的代理服务器实例
func New(opts ...ServerOption) *ProxyServe {
	ctx, cancel := context.WithCancel(context.Background())

	s := &ProxyServe{
		config:         &DefaultConfig,
		logger:         GetLogger(),
		maxConnections: 10000, // 默认最大连接数
		shutdownCtx:    ctx,
		shutdownCancel: cancel,
		signalChan:     make(chan os.Signal, 1),
	}

	// 应用选项
	for _, opt := range opts {
		opt(s)
	}

	// 初始化上下文池
	s.contextPool.New = func() interface{} {
		return s.newContext()
	}

	return s
}

// WithConfig 使用配置
func WithConfig(cfg *ProxyConfig) ServerOption {
	return func(s *ProxyServe) {
		if cfg != nil {
			s.config = cfg
		}
	}
}

// WithLogger 使用自定义日志器
func WithLogger(logger *Logger) ServerOption {
	return func(s *ProxyServe) {
		if logger != nil {
			s.logger = logger
		}
	}
}

// WithMaxConnections 设置最大连接数
func WithMaxConnections(max int64) ServerOption {
	return func(s *ProxyServe) {
		s.maxConnections = max
	}
}

// WithMethod 设置代理方法
func WithMethod(method string) ServerOption {
	return func(s *ProxyServe) {
		s.config.Method = method
	}
}

// WithHost 设置监听地址
func WithHost(ip, port string) ServerOption {
	return func(s *ProxyServe) {
		s.config.IP = ip
		s.config.Port = port
	}
}

// ========== 链式调用 API ==========

// Use 添加中间件（链式调用）
func (s *ProxyServe) Use(middleware ...ProxyHandle) *ProxyServe {
	s.handlersChain = append(s.handlersChain, middleware...)
	return s
}

// Handle 添加处理器（链式调用）
func (s *ProxyServe) Handle(handler ProxyHandle) *ProxyServe {
	s.handlersChain = append(s.handlersChain, handler)
	return s
}

// SetMethod 设置代理方法（链式调用）
func (s *ProxyServe) SetMethod(method string) *ProxyServe {
	s.config.Method = method
	return s
}

// SetHost 设置监听地址（链式调用）
func (s *ProxyServe) SetHost(ip, port string) *ProxyServe {
	s.config.IP = ip
	s.config.Port = port
	return s
}

// SetMaxConnections 设置最大连接数（链式调用）
func (s *ProxyServe) SetMaxConnections(max int64) *ProxyServe {
	s.maxConnections = max
	return s
}

// ========== 向后兼容 API ==========

// AddHandle 设置自定义代理响应处理（从尾部添加）
// Deprecated: 使用 Handle() 代替
func (s *ProxyServe) AddHandle(proxyHandle ProxyHandle) {
	s.handlersChain = append(s.handlersChain, proxyHandle)
}

// AddMiddleware 从头部添加中间件
// Deprecated: 使用 Use() 代替
func (s *ProxyServe) AddMiddleware(proxyHandle ProxyHandle) {
	s.handlersChain = append([]ProxyHandle{proxyHandle}, s.handlersChain...)
}

// ========== 服务启动与管理 ==========

// Start 启动服务器
func (s *ProxyServe) Start() error {
	return s.StartWithContext(context.Background())
}

// StartWithContext 使用上下文启动服务器
func (s *ProxyServe) StartWithContext(ctx context.Context) error {
	// 向后兼容：处理旧字段
	s.migrateLegacyFields()

	// 验证配置
	if err := s.config.Validate(); err != nil {
		return NewProxyError("validate", "", err)
	}

	// 初始化处理器链
	if err := s.initHandlers(); err != nil {
		return err
	}

	// 初始化日志
	s.initLogger()

	// 启动监听
	addr := s.config.IP + ":" + s.config.Port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return NewProxyError("listen", addr, err)
	}
	s.listener = ln
	s.running = true
	s.stats.StartTime = time.Now()

	s.logger.Info("Server started on %s, method: %s", addr, s.config.Method)
	s.printBanner()

	// 启动信号处理协程
	go s.handleSignals()

	// 接受连接循环
	return s.acceptLoop(ctx)
}

// Run 启动服务器并处理优雅关闭（便捷方法）
func (s *ProxyServe) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// 等待关闭信号
	<-s.shutdownCtx.Done()
	return nil
}

// ServerHandleListen ServerHandle 服务代理处理
// Deprecated: 使用 Start() 代替
func (s *ProxyServe) ServerHandleListen() {
	ServeProxy(s)
}

// migrateLegacyFields 迁移旧字段到新配置
func (s *ProxyServe) migrateLegacyFields() {
	// 确保 config 已初始化
	if s.config == nil {
		s.config = &DefaultConfig
	}

	// 确保 logger 已初始化
	if s.logger == nil {
		s.logger = GetLogger()
	}

	// 确保上下文池已初始化
	if s.contextPool.New == nil {
		s.contextPool.New = func() interface{} {
			return &Context{
				Handlers: s.handlersChain,
			}
		}
	}

	// 确保信号通道已初始化
	if s.signalChan == nil {
		ctx, cancel := context.WithCancel(context.Background())
		s.shutdownCtx = ctx
		s.shutdownCancel = cancel
		s.signalChan = make(chan os.Signal, 1)
	}

	if s.Ip != "" {
		s.config.IP = s.Ip
	}
	if s.Port != "" {
		s.config.Port = s.Port
	}
	if s.Method != "" {
		s.config.Method = s.Method
	}
	if s.Protocol != "" {
		s.config.Protocol = s.Protocol
	}
	if s.LogPath != "" {
		s.config.LogPath = s.LogPath
	}
	// 处理旧的处理器链
	if len(s.ProxyHandlersChain) > 0 && len(s.handlersChain) == 0 {
		s.handlersChain = s.ProxyHandlersChain
	}
}

// initHandlers 初始化处理器链
func (s *ProxyServe) initHandlers() error {
	// 根据方法添加默认处理器
	switch s.config.Method {
	case HttpProxy:
		s.handlersChain = append(s.handlersChain, HttpProxyHandle)
	case HttpsProxy:
		s.handlersChain = append(s.handlersChain, HttpsProxyHandle)
	case HttpsSniffProxy:
		s.handlersChain = append(s.handlersChain, HttpsSniffProxyHandle)
	case HttpsSniffDetailProxy:
		s.handlersChain = append(s.handlersChain, HttpsSniffDetailProxyHandle)
	case HttpInterceptProxy:
		s.handlersChain = append(s.handlersChain, HttpInterceptProxyHandle)
	case TcpTunnelProxy:
		s.handlersChain = append(s.handlersChain, TunnelProxyHandle)
	case WebsocketProxy:
		s.handlersChain = append(s.handlersChain, WebsocketProxyHandle)
	case DefaultProxy:
		if len(s.handlersChain) == 0 {
			return NewProxyError("init", "", ErrInvalidRequest)
		}
	default:
		s.logger.Warn("Unknown method: %s, using default", s.config.Method)
	}

	return nil
}

// initLogger 初始化日志
func (s *ProxyServe) initLogger() {
	if s.config.LogPath != "" {
		s.logger.SetOutput(NewFileWriter(s.config.LogPath))
	}
}

// acceptLoop 接受连接循环
func (s *ProxyServe) acceptLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.shutdownCtx.Done():
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				// 检查是否是关闭导致的错误
				if s.shuttingDown {
					return nil
				}

				// 记录错误但继续接受新连接
				s.logger.Error("Accept error: %v", err)
				atomic.AddInt64(&s.stats.FailedConnections, 1)
				continue
			}

			// 检查连接数限制
			if atomic.LoadInt64(&s.stats.ActiveConnections) >= s.maxConnections {
				s.logger.Warn("Max connections reached, rejecting: %s", conn.RemoteAddr())
				conn.Close()
				continue
			}

			// 处理连接
			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理单个连接
func (s *ProxyServe) handleConnection(conn net.Conn) {
	// 更新统计
	atomic.AddInt64(&s.stats.TotalConnections, 1)
	atomic.AddInt64(&s.stats.ActiveConnections, 1)
	s.connWg.Add(1)

	defer func() {
		atomic.AddInt64(&s.stats.ActiveConnections, -1)
		s.connWg.Done()
		conn.Close()
	}()

	// 获取上下文
	ctx := s.contextPool.Get().(*Context)
	defer func() {
		ctx.Reset()
		s.contextPool.Put(ctx)
	}()

	// 设置客户端连接
	ctx.ClientConn = conn
	ctx.StartTime = time.Now()

	// 执行处理器链
	ctx.Next()
}

// Shutdown 优雅关闭服务器
func (s *ProxyServe) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	if s.shuttingDown {
		s.mu.Unlock()
		return nil
	}
	s.shuttingDown = true
	s.mu.Unlock()

	s.logger.Info("Shutting down server...")

	// 关闭监听器
	if s.listener != nil {
		s.listener.Close()
	}

	// 取消上下文
	s.shutdownCancel()

	// 等待所有连接处理完成
	done := make(chan struct{})
	go func() {
		s.connWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("Server shutdown completed")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Shutdown timeout, forcing close")
		return ctx.Err()
	}
}

// handleSignals 处理系统信号
func (s *ProxyServe) handleSignals() {
	signal.Notify(s.signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-s.signalChan:
		s.logger.Info("Received signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	case <-s.shutdownCtx.Done():
		return
	}
}

// Stats 获取服务器统计信息
func (s *ProxyServe) Stats() ServerStats {
	return ServerStats{
		TotalConnections:  atomic.LoadInt64(&s.stats.TotalConnections),
		ActiveConnections: atomic.LoadInt64(&s.stats.ActiveConnections),
		FailedConnections: atomic.LoadInt64(&s.stats.FailedConnections),
		BytesReceived:     atomic.LoadInt64(&s.stats.BytesReceived),
		BytesSent:         atomic.LoadInt64(&s.stats.BytesSent),
		StartTime:         s.stats.StartTime,
	}
}

// IsRunning 检查服务器是否运行中
func (s *ProxyServe) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running && !s.shuttingDown
}

// newContext 创建新上下文
func (s *ProxyServe) newContext() *Context {
	return &Context{
		Handlers: s.handlersChain,
	}
}

// printBanner 打印启动横幅
func (s *ProxyServe) printBanner() {
	println("/\\  _`\\    __/\\  _`\\                                \n\\ \\ \\/\\_\\ /\\_\\ \\ \\L\\ \\_ __   ___   __  _  __  __    \n \\ \\ \\/_/_\\/\\ \\ \\ ,__/\\`'__\\/ __`\\/\\ \\/'\\/\\ \\/\\ \\   \n  \\ \\ \\L\\ \\\\ \\ \\ \\ \\/\\ \\ \\//\\ \\L\\ \\/>  </\\ \\ \\_\\ \\  \n   \\ \\____/ \\ \\_\\ \\_\\ \\ \\_\\\\ \\____//\\_/\\_\\\\/`____ \\ \n    \\/___/   \\/_/\\/_/  \\/_/ \\/___/ \\//\\/_/ `/___/> \\\n                                              /\\___/\n                                              \\/__/ ")
	fmt.Printf("CiProxy Version %s, Mode %s\n", ProxyVersion, ProxyMode)
	s.logger.Info("Listen on %s:%s, Proxy Method %s, Protocol %s", s.config.IP, s.config.Port, s.config.Method, s.config.Protocol)
}

// FileWriter 文件写入器
type FileWriter struct {
	file *os.File
}

// NewFileWriter 创建文件写入器
func NewFileWriter(path string) *FileWriter {
	// 确保目录存在
	os.MkdirAll("./log/", os.ModePerm)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return &FileWriter{file: os.Stdout}
	}
	return &FileWriter{file: file}
}

// Write 实现 io.Writer 接口
func (w *FileWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}