/**
 * @file mitm_proxy_server/main.go
 * @brief MITM拦截代理服务器示例
 * @description 展示如何使用CiProxy创建完整的MITM拦截代理，支持请求/响应修改
 *
 * 功能特点:
 *   - 支持HTTP/1.1 和 HTTP/2 协议拦截
 *   - 支持请求/响应拦截、修改、阻断
 *   - 内置流量捕获功能
 *   - 支持中间件扩展
 *
 * 使用方法:
 *   基本模式:   go run main.go -ip 127.0.0.1 -port 8080
 *   流量捕获:   go run main.go -ip 127.0.0.1 -port 8080 -capture
 *   禁用HTTP/2: go run main.go -ip 127.0.0.1 -port 8080 -http2=false
 *
 * 测试方法:
 *   curl -x http://127.0.0.1:8080 https://httpbin.org/get -k
 *
 * @author cilang
 */

package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/opencvlzg/ciproxy"
	"github.com/opencvlzg/ciproxy/pkg/middleware/builtins"
	mitm "github.com/opencvlzg/ciproxy/pkg/module/mitm"
	"github.com/opencvlzg/ciproxy/pkg/transfer"
)

func main() {
	// 命令行参数解析
	ip := flag.String("ip", "127.0.0.1", "监听IP地址")
	port := flag.String("port", "8080", "监听端口")
	enableCapture := flag.Bool("capture", false, "启用流量捕获")
	enableHTTP2 := flag.Bool("http2", true, "启用HTTP/2支持")
	flag.Parse()

	// 配置MITM拦截器
	interceptorConfig := &mitm.InterceptorConfig{
		EnableTrafficCapture: *enableCapture,
		EnableHTTP2:          *enableHTTP2,
	}
	ciproxy.SetInterceptorConfig(interceptorConfig)

	// 获取拦截器实例
	interceptor := ciproxy.GetInterceptor()

	// ========== 添加内置中间件 ==========

	// 1. 日志中间件
	interceptor.Use(&builtins.LoggingMiddleware{
		Logger: log.Printf,
	})

	// 2. 请求头修改中间件
	interceptor.Use(&builtins.HeaderModifier{
		RequestHeaders: map[string]string{
			"X-Proxy-By": "CiProxy-MITM",
		},
	})

	// ========== 添加自定义拦截中间件 ==========

	// 3. 请求阻断中间件 - 阻断特定域名
	interceptor.UseFunc(func(c *ciproxy.Context) *http.Response {
		if c.GetRequest() == nil {
			return nil
		}

		url := c.GetRequest().URL.String()

		// 阻断广告和追踪域名
		if isBlockedURL(url) {
			log.Printf("[BLOCK] 阻断请求: %s", url)
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"error": "blocked by proxy"}`)),
			}
		}

		// 修改请求头示例
		c.SetRequestHeader("X-Intercepted-By", "CiProxy")

		log.Printf("[INTERCEPT] %s %s", c.GetRequest().Method, url)
		return nil // 返回nil继续处理
	})

	// 创建MITM代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.HttpInterceptProxy)

	// 启动服务器
	go func() {
		log.Printf("MITM拦截代理启动: %s:%s", *ip, *port)
		log.Printf("流量捕获: %v, HTTP/2支持: %v", *enableCapture, *enableHTTP2)
		log.Println("")
		log.Println("已启用的功能:")
		log.Println("  1. 日志记录中间件")
		log.Println("  2. 请求头修改中间件")
		log.Println("  3. 请求阻断中间件")
		log.Println("")
		log.Println("测试方法: curl -x http://127.0.0.1:8080 https://httpbin.org/get -k")

		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 如果启用了流量捕获，定期打印统计信息
	if *enableCapture {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				store := transfer.GetTrafficStore()
				log.Printf("[STATS] 捕获流量: %d 条", store.Len())
			}
		}()
	}

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 如果启用了流量捕获，导出数据
	if *enableCapture {
		exportTraffic()
	}

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}

// isBlockedURL 判断是否需要阻断的URL
func isBlockedURL(url string) bool {
	// 示例：阻断广告和追踪域名
	blockedDomains := []string{
		"ads.example.com",
		"tracking.example.com",
		"analytics.example.com",
	}

	for _, domain := range blockedDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}

// exportTraffic 导出流量数据
func exportTraffic() {
	store := transfer.GetTrafficStore()
	if store.Len() == 0 {
		log.Println("没有捕获到流量数据")
		return
	}

	data, err := store.Export()
	if err != nil {
		log.Printf("导出流量失败: %v", err)
		return
	}

	filename := "captured_traffic_" + time.Now().Format("20060102_150405") + ".json"
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("写入文件失败: %v", err)
		return
	}

	log.Printf("流量数据已导出到: %s (%d 条)", filename, store.Len())
}