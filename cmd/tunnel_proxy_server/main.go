/**
 * @file tunnel_proxy_server/main.go
 * @brief TCP隧道代理服务器示例
 * @description 展示如何使用CiProxy创建TCP隧道代理
 *
 * 功能特点:
 *   - 支持TCP协议隧道
 *   - 高性能流量转发
 *   - 支持大并发连接
 *   - 支持流量加密（可选）
 *
 * 使用方法:
 *   基本模式:   go run main.go -ip 127.0.0.1 -port 8080
 *   高并发模式: go run main.go -ip 0.0.0.0 -port 8080 -max-conn 50000
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:8080 https://httpbin.org/ip
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

func main() {
	// 命令行参数解析
	ip := flag.String("ip", "127.0.0.1", "监听IP地址")
	port := flag.String("port", "8080", "监听端口")
	maxConn := flag.Int64("max-conn", 50000, "最大并发连接数")
	flag.Parse()

	// 创建TCP隧道代理
	// 使用选项模式和链式调用混合风格
	server := ciproxy.New(
		ciproxy.WithHost(*ip, *port),
		ciproxy.WithMethod(ciproxy.TcpTunnelProxy),
		ciproxy.WithMaxConnections(*maxConn),
	).Use(func(c *ciproxy.Context) {
		// 简单的连接日志
		log.Printf("[TUNNEL] %s -> forwarding", c.ClientConn.RemoteAddr())
		c.Next()
	})

	// 启动服务器
	go func() {
		log.Printf("TCP隧道代理启动: %s:%s", *ip, *port)
		log.Printf("最大连接数: %d", *maxConn)
		log.Println("隧道模式: 高性能流量转发")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 获取服务器统计信息
	stats := server.Stats()
	log.Printf("服务器统计: 总连接=%d, 活跃连接=%d",
		stats.TotalConnections, stats.ActiveConnections)

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}