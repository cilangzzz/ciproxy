/**
 * @file websocket_proxy_server/main.go
 * @brief WebSocket代理服务器示例
 * @description 展示如何使用CiProxy创建WebSocket代理
 *
 * 功能特点:
 *   - 支持WebSocket协议代理
 *   - 自动处理协议升级
 *   - 支持长连接
 *   - 支持ws和wss协议
 *
 * 使用方法:
 *   go run main.go -ip 127.0.0.1 -port 8080
 *
 * 测试方法:
 *   配置浏览器或WebSocket客户端使用此代理
 *   例如: wss://echo.websocket.org 通过代理连接
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
	maxConn := flag.Int64("max-conn", 10000, "最大连接数")
	flag.Parse()

	// 创建WebSocket代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.WebsocketProxy).
		SetMaxConnections(*maxConn).
		Use(func(c *ciproxy.Context) {
			// 记录WebSocket连接信息
			if c.GetRequest() != nil {
				req := c.GetRequest()
				log.Printf("[WS] %s %s", req.Method, req.Host)

				// 检查是否是WebSocket升级请求
				if req.Header.Get("Upgrade") == "websocket" {
					log.Printf("[WS] WebSocket升级请求: %s", req.URL.String())
				}
			}

			c.Next()

			// 连接结束后的日志
			log.Printf("[WS] 连接结束: %s", c.ConnStatus)
		})

	// 启动服务器
	go func() {
		log.Printf("WebSocket代理启动: %s:%s", *ip, *port)
		log.Printf("最大连接数: %d", *maxConn)
		log.Println("支持协议: ws, wss")
		log.Println("提示: 配置WebSocket客户端使用此代理服务器")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}