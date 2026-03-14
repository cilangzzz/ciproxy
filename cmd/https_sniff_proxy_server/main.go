/**
 * @file https_sniff_proxy_server/main.go
 * @brief HTTPS嗅探代理服务器示例
 * @description 展示如何使用CiProxy创建HTTPS嗅探代理，解密并查看HTTPS流量
 *
 * 功能特点:
 *   - 支持HTTPS流量解密嗅探
 *   - 动态生成TLS证书
 *   - 支持详细模式和普通模式
 *   - 可查看请求/响应内容
 *
 * 注意事项:
 *   - 需要将CA证书安装到系统/浏览器信任列表
 *   - 默认使用内置的测试证书（仅供开发测试）
 *   - 生产环境请使用正式CA证书
 *
 * 使用方法:
 *   普通模式:   go run main.go -ip 127.0.0.1 -port 6677 -mode sniff
 *   详细模式:   go run main.go -ip 127.0.0.1 -port 6677 -mode detail
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:6677 https://httpbin.org/ip -k
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
	mode := flag.String("mode", "sniff", "嗅探模式: sniff(普通) 或 detail(详细)")
	flag.Parse()

	// 根据模式选择代理方法
	var method string
	switch *mode {
	case "detail":
		method = ciproxy.HttpsSniffDetailProxy
		log.Println("使用详细嗅探模式 (HttpsSniffDetailProxy)")
	default:
		method = ciproxy.HttpsSniffProxy
		log.Println("使用普通嗅探模式 (HttpsSniffProxy)")
	}

	// 创建嗅探代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(method).
		Use(func(c *ciproxy.Context) {
			// 记录TLS连接信息
			if c.TlsClientConn != nil && c.TlsServerConn != nil {
				log.Printf("[SNIFF] 客户端: %s -> 服务端: %s",
					c.TlsClientConn.RemoteAddr(),
					c.TlsServerConn.RemoteAddr())
			}

			// 记录请求信息（如果在详细模式下）
			if c.GetRequest() != nil {
				req := c.GetRequest()
				log.Printf("[REQUEST] %s %s", req.Method, req.URL.String())
			}

			c.Next()
		})

	// 启动服务器
	go func() {
		log.Printf("HTTPS嗅探代理启动: %s:%s", *ip, *port)
		log.Println("提示: 请确保已将CA证书安装到系统信任列表")
		log.Println("证书位置: ./cert/root.crt")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}