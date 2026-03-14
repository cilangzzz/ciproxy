/**
 * @file http_proxy_server/main.go
 * @brief HTTP代理服务器示例
 * @description 展示如何使用CiProxy创建基本的HTTP代理服务器
 *
 * 功能特点:
 *   - 支持HTTP协议代理
 *   - 链式调用API风格
 *   - 内置日志记录
 *   - 优雅关闭支持
 *
 * 使用方法:
 *   go run main.go -ip 127.0.0.1 -port 8080
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:8080 http://httpbin.org/ip
 *
 * @author cilang
 */

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opencvlzg/ciproxy"
)

// loggingMiddleware 自定义日志中间件
func loggingMiddleware(c *ciproxy.Context) {
	start := time.Now()

	// 记录请求开始
	log.Printf("[REQUEST] %s -> %s",
		c.ClientConn.RemoteAddr(),
		c.ConnStatus)

	// 执行下一个处理器
	c.Next()

	// 记录请求完成
	duration := time.Since(start)
	log.Printf("[RESPONSE] 耗时: %v", duration)
}

func main() {
	// 命令行参数解析
	ip := flag.String("ip", "127.0.0.1", "监听IP地址")
	port := flag.String("port", "8080", "监听端口")
	flag.Parse()

	// 使用链式调用API创建HTTP代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.HttpProxy).
		Use(loggingMiddleware)

	// 启动服务器（在后台运行）
	go func() {
		log.Printf("HTTP代理服务器启动: %s:%s", *ip, *port)
		log.Println("使用方法: curl -x http://127.0.0.1:8080 http://httpbin.org/ip")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭错误: %v", err)
	}
	log.Println("服务器已关闭")
}