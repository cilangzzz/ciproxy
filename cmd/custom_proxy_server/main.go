/**
 * @file custom_proxy_server/main.go
 * @brief 自定义代理服务器示例
 * @description 展示如何使用自定义处理器实现完全自定义的代理逻辑
 *
 * 功能特点:
 *   - 完全自定义的请求处理逻辑
 *   - 支持自定义协议处理
 *   - 支持请求/响应修改
 *   - 支持条件性阻断
 *
 * 使用方法:
 *   go run main.go -ip 127.0.0.1 -port 8888
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:8888 http://httpbin.org/get
 *
 * @author cilang
 */

package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/opencvlzg/ciproxy"
)

// CustomProxyHandler 自定义代理处理器
// 实现完全自定义的代理逻辑
func CustomProxyHandler(c *ciproxy.Context) {
	// 1. 读取客户端请求
	buf := bufio.NewReader(c.ClientConn)
	request, err := http.ReadRequest(buf)
	if err != nil {
		log.Printf("[ERROR] 读取请求失败: %v", err)
		return
	}

	// 保存请求到上下文
	c.SetRequest(request)

	log.Printf("[REQUEST] %s %s", request.Method, request.Host)

	// 2. 自定义请求处理逻辑

	// 示例：修改请求头
	request.Header.Set("X-Custom-Proxy", "CiProxy-Custom")
	request.Header.Del("X-Forwarded-For")

	// 示例：根据条件阻断请求
	if shouldBlockRequest(request) {
		log.Printf("[BLOCK] 阻断请求: %s", request.Host)
		c.BlockWithStatus(http.StatusForbidden, "Access Denied by Custom Proxy")
		return
	}

	// 3. 建立到目标服务器的连接
	host := request.Host
	if !strings.Contains(host, ":") {
		if request.Method == "CONNECT" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	serverConn, err := net.DialTimeout("tcp", host, ciproxy.DefaultOutTime)
	if err != nil {
		log.Printf("[ERROR] 连接服务器失败: %s, %v", host, err)
		c.BlockWithStatus(http.StatusBadGateway, "Bad Gateway")
		return
	}
	defer serverConn.Close()

	c.SetServerConn(serverConn)

	// 4. 处理HTTPS CONNECT请求
	if request.Method == "CONNECT" {
		_, err = c.ClientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			log.Printf("[ERROR] 发送响应失败: %v", err)
			return
		}
		log.Printf("[CONNECT] 隧道建立: %s", host)
	} else {
		// 对于HTTP请求，转发修改后的请求
		if err := request.Write(serverConn); err != nil {
			log.Printf("[ERROR] 转发请求失败: %v", err)
			return
		}
	}

	// 5. 双向数据转发
	go ciproxy.Transfer(c.ClientConn, serverConn)
	go ciproxy.Transfer(serverConn, c.ClientConn)

	// 6. 更新连接状态
	c.ConnStatus = "connected"
}

// shouldBlockRequest 判断是否应该阻断请求
func shouldBlockRequest(req *http.Request) bool {
	// 示例：阻断特定域名
	blockedHosts := []string{
		"ads.example.com",
		"tracking.example.com",
	}

	for _, blocked := range blockedHosts {
		if strings.Contains(req.Host, blocked) {
			return true
		}
	}

	return false
}

// RequestLoggerMiddleware 请求日志中间件
func RequestLoggerMiddleware(c *ciproxy.Context) {
	start := time.Now()

	// 记录请求开始
	log.Printf("[START] %s", c.ClientConn.RemoteAddr())

	// 执行下一个处理器
	c.Next()

	// 记录请求完成
	duration := time.Since(start)
	log.Printf("[END] 耗时: %v, 状态: %s", duration, c.ConnStatus)
}

func main() {
	// 命令行参数解析
	ip := flag.String("ip", "127.0.0.1", "监听IP地址")
	port := flag.String("port", "8888", "监听端口")
	flag.Parse()

	// 创建自定义代理服务器
	// 使用 DefaultProxy 模式，需要自行实现所有处理逻辑
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.DefaultProxy).
		Use(RequestLoggerMiddleware).   // 添加中间件
		Handle(CustomProxyHandler)        // 添加自定义处理器

	// 启动服务器
	go func() {
		log.Printf("自定义代理服务器启动: %s:%s", *ip, *port)
		log.Println("使用自定义处理器: CustomProxyHandler")
		log.Println("")
		log.Println("功能特点:")
		log.Println("  - 自定义请求头修改")
		log.Println("  - 条件性请求阻断")
		log.Println("  - 请求日志记录")
		log.Println("")
		log.Println("测试方法: curl -x http://127.0.0.1:8888 http://httpbin.org/get")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}