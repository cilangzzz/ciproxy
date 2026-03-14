/**
 * @file middleware_demo/main.go
 * @brief 中间件使用示例
 * @description 展示如何使用内置中间件和自定义中间件
 *
 * 功能特点:
 *   - 展示所有内置中间件的使用
 *   - 展示如何编写自定义中间件
 *   - 展示中间件执行顺序
 *   - 展示CORS跨域配置
 *
 * 使用方法:
 *   go run main.go -ip 127.0.0.1 -port 8080
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
)

func main() {
	// 命令行参数解析
	ip := flag.String("ip", "127.0.0.1", "监听IP地址")
	port := flag.String("port", "8080", "监听端口")
	flag.Parse()

	// 配置MITM拦截器
	ciproxy.SetInterceptorConfig(&mitm.InterceptorConfig{
		EnableTrafficCapture: false,
		EnableHTTP2:          true,
	})

	interceptor := ciproxy.GetInterceptor()

	// ========== 1. 日志中间件 ==========
	interceptor.Use(&builtins.LoggingMiddleware{
		Logger: log.Printf,
	})

	// ========== 2. 请求头修改中间件 ==========
	interceptor.Use(&builtins.HeaderModifier{
		RequestHeaders: map[string]string{
			"X-Proxy-Version": "1.0",
			"X-Custom-Header": "CiProxy-Demo",
		},
		DelRequestHeaders: []string{
			"X-Forwarded-For",
			"Via",
		},
		ResponseHeaders: map[string]string{
			"X-Served-By": "CiProxy",
		},
	})

	// ========== 3. CORS跨域中间件 ==========
	interceptor.Use(&builtins.CorsMiddleware{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// ========== 4. 请求阻断中间件 ==========
	interceptor.Use(&builtins.BlockMiddleware{
		URLPatterns: []string{
			`.*\.ads\.example\.com.*`, // 阻断广告域名
			`.*\/tracking\/.*`,         // 阻断追踪路径
		},
		StatusCode:   http.StatusForbidden,
		ResponseBody: `{"error": "blocked by proxy"}`,
		ResponseHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	})

	// ========== 5. 自定义中间件示例 ==========

	// 5.1 认证检查中间件
	interceptor.UseFunc(func(c *ciproxy.Context) *http.Response {
		req := c.GetRequest()
		if req == nil {
			return nil
		}

		// 检查特定路径的认证
		if strings.HasPrefix(req.URL.Path, "/admin/") {
			auth := req.Header.Get("Authorization")
			if auth == "" {
				log.Printf("[AUTH] 未授权访问: %s", req.URL.Path)
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(`{"error": "unauthorized"}`)),
				}
			}
		}

		return nil
	})

	// 5.2 计时中间件
	interceptor.UseFunc(func(c *ciproxy.Context) *http.Response {
		start := time.Now()
		log.Printf("[TIMING] 请求开始: %s", start.Format(time.RFC3339Nano))
		// 注意：UseFunc 不支持响应后处理，仅记录开始时间
		return nil
	})

	// 5.3 IP白名单中间件
	interceptor.UseFunc(func(c *ciproxy.Context) *http.Response {
		// 获取客户端IP
		clientAddr := c.ClientConn.RemoteAddr().String()

		// 检查IP白名单（示例：只允许本地访问）
		whitelist := []string{"127.0.0.1", "::1", "localhost"}
		allowed := false
		for _, allowedIP := range whitelist {
			if strings.Contains(clientAddr, allowedIP) {
				allowed = true
				break
			}
		}

		if !allowed {
			log.Printf("[IP-FILTER] 拒绝访问: %s", clientAddr)
			return &http.Response{
				StatusCode: http.StatusForbidden,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"error": "IP not allowed"}`)),
			}
		}

		log.Printf("[IP-FILTER] 允许访问: %s", clientAddr)
		return nil
	})

	// 创建代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.HttpInterceptProxy)

	// 启动服务器
	go func() {
		log.Printf("中间件演示服务器启动: %s:%s", *ip, *port)
		log.Println("")
		log.Println("已启用的中间件:")
		log.Println("  1. LoggingMiddleware - 日志记录")
		log.Println("  2. HeaderModifier - 请求头修改")
		log.Println("  3. CorsMiddleware - CORS跨域")
		log.Println("  4. BlockMiddleware - 请求阻断")
		log.Println("  5. 自定义中间件 - 认证、计时、IP过滤")
		log.Println("")
		log.Println("测试方法: curl -x http://127.0.0.1:8080 https://httpbin.org/get -k")

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