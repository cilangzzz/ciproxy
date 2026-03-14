/**
 * @file tunnel_proxy_client/main.go
 * @brief TCP隧道代理客户端示例
 * @description 展示如何创建TCP隧道代理客户端，连接远程代理服务器
 *
 * 功能特点:
 *   - 连接远程代理服务器
 *   - 本地代理服务
 *   - 支持加密隧道
 *
 * 使用方法:
 *   go run main.go -server 127.0.0.1:8080 -local 127.0.0.1:8888
 *
 * @author cilang
 */

package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opencvlzg/ciproxy"
)

func main() {
	// 命令行参数解析
	serverAddr := flag.String("server", "127.0.0.1:8080", "远程代理服务器地址")
	localIP := flag.String("local-ip", "127.0.0.1", "本地监听IP")
	localPort := flag.String("local-port", "8888", "本地监听端口")
	flag.Parse()

	log.Printf("TCP隧道代理客户端启动")
	log.Printf("远程服务器: %s", *serverAddr)
	log.Printf("本地监听: %s:%s", *localIP, *localPort)

	// 启动本地代理服务器
	server := ciproxy.New().
		SetHost(*localIP, *localPort).
		SetMethod(ciproxy.TcpTunnelProxy).
		Use(func(c *ciproxy.Context) {
			// 记录连接信息
			log.Printf("[TUNNEL] 本地连接: %s -> 远程服务器: %s",
				c.ClientConn.RemoteAddr(), *serverAddr)
			c.Next()
		})

	// 测试与远程服务器的连接
	conn, err := net.DialTimeout("tcp", *serverAddr, 5*time.Second)
	if err != nil {
		log.Fatalf("无法连接到远程服务器 %s: %v", *serverAddr, err)
	}
	conn.Close()
	log.Printf("成功连接到远程服务器: %s", *serverAddr)

	// 启动服务器
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	log.Printf("本地代理服务已启动: %s:%s", *localIP, *localPort)
	log.Println("提示: 按 Ctrl+C 停止服务器")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("服务器已关闭")
}