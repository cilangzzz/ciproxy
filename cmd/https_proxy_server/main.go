/**
 * @file https_proxy_server/main.go
 * @brief HTTPS隧道代理服务器示例
 * @description 展示如何使用CiProxy创建HTTPS隧道代理服务器
 *
 * 功能特点:
 *   - 支持HTTPS协议隧道代理
 *   - 自动处理CONNECT方法
 *   - 支持链式调用和选项模式两种API风格
 *   - 支持优雅关闭
 *
 * 使用方法:
 *   方式1: go run main.go -ip 127.0.0.1 -port 6677
 *   方式2: go run main.go -ip 0.0.0.0 -port 443
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:6677 https://httpbin.org/ip
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
	port := flag.String("port", "6677", "监听端口")
	maxConn := flag.Int64("max-conn", 10000, "最大连接数")
	flag.Parse()

	// 方式1: 链式调用API (推荐)
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.HttpsProxy).
		SetMaxConnections(*maxConn).
		Use(func(c *ciproxy.Context) {
			// 简单的连接日志中间件
			log.Printf("[CONNECT] %s -> %s",
				c.ClientConn.RemoteAddr(),
				c.ConnStatus)
			c.Next()
		})

	// 方式2: 使用选项模式 (等价写法)
	// server := ciproxy.New(
	//     ciproxy.WithHost(*ip, *port),
	//     ciproxy.WithMethod(ciproxy.HttpsProxy),
	//     ciproxy.WithMaxConnections(*maxConn),
	// )

	// 启动服务器 (在后台运行)
	go func() {
		log.Printf("HTTPS代理服务器启动中: %s:%s", *ip, *port)
		log.Printf("最大连接数: %d", *maxConn)
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待中断信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 给予30秒时间处理未完成的连接
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭错误: %v", err)
	}
	log.Println("服务器已关闭")
}