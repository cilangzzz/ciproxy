/**
 * @file traffic_capture_server/main.go
 * @brief 流量捕获代理服务器示例
 * @description 展示如何使用CiProxy捕获和存储HTTP/HTTPS流量
 *
 * 功能特点:
 *   - 自动捕获所有请求/响应
 *   - 支持流量过滤和搜索
 *   - 支持流量导出为JSON格式
 *   - 支持流量统计信息
 *
 * 使用方法:
 *   go run main.go -ip 127.0.0.1 -port 8080 -output traffic.json
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
	"fmt"
	"log"
	"os"
	"os/signal"
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
	output := flag.String("output", "", "流量输出文件路径 (默认自动生成)")
	maxSize := flag.Int("max-size", 10000, "最大存储条目数")
	printStats := flag.Bool("stats", true, "定期打印统计信息")
	flag.Parse()

	// 配置MITM拦截器，启用流量捕获
	ciproxy.SetInterceptorConfig(&mitm.InterceptorConfig{
		EnableTrafficCapture: true,
		EnableHTTP2:          true,
	})

	interceptor := ciproxy.GetInterceptor()

	// 添加日志中间件
	interceptor.Use(&builtins.LoggingMiddleware{
		Logger: func(format string, args ...interface{}) {
			log.Printf("[TRAFFIC] "+format, args...)
		},
	})

	// 创建代理服务器
	server := ciproxy.New().
		SetHost(*ip, *port).
		SetMethod(ciproxy.HttpInterceptProxy)

	// 获取流量存储实例
	store := transfer.GetTrafficStore()
	store.SetMaxEntries(*maxSize)

	// 定期打印统计信息
	if *printStats {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				stats := store.Stats()
				log.Printf("[STATS] 总条目: %d/%d, 平均耗时: %v",
					stats.TotalEntries, stats.MaxEntries, stats.AvgDuration)
			}
		}()
	}

	// 启动服务器
	go func() {
		log.Printf("流量捕获代理启动: %s:%s", *ip, *port)
		log.Printf("最大存储条目: %d", *maxSize)
		log.Println("提示: 按 Ctrl+C 停止服务器并导出流量数据")
		log.Println("")
		log.Println("测试方法: curl -x http://127.0.0.1:8080 https://httpbin.org/get -k")
		if err := server.Run(); err != nil {
			log.Fatalf("服务器错误: %v", err)
		}
	}()

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 导出流量数据
	exportTraffic(store, *output)

	// 打印最终统计
	printFinalStats(store)

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("服务器已关闭")
}

// exportTraffic 导出流量数据
func exportTraffic(store *transfer.TrafficStore, outputPath string) {
	if store.Len() == 0 {
		log.Println("没有捕获到流量数据")
		return
	}

	// 生成输出文件名
	if outputPath == "" {
		outputPath = fmt.Sprintf("traffic_%s.json", time.Now().Format("20060102_150405"))
	}

	data, err := store.Export()
	if err != nil {
		log.Printf("导出失败: %v", err)
		return
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Printf("写入文件失败: %v", err)
		return
	}
	log.Printf("流量已导出到: %s (%d 条)", outputPath, store.Len())
}

// printFinalStats 打印最终统计信息
func printFinalStats(store *transfer.TrafficStore) {
	stats := store.Stats()
	entries := store.List()

	log.Println("========== 流量统计 ==========")
	log.Printf("总捕获条目: %d", stats.TotalEntries)
	log.Printf("平均请求耗时: %v", stats.AvgDuration)

	// 按方法统计
	methodCount := make(map[string]int)
	statusCount := make(map[int]int)
	hostCount := make(map[string]int)

	for _, e := range entries {
		methodCount[e.Request.Method]++
		statusCount[e.Response.Status]++
		hostCount[e.Request.Host]++
	}

	log.Println("\n请求方法分布:")
	for method, count := range methodCount {
		log.Printf("  %s: %d", method, count)
	}

	log.Println("\n状态码分布:")
	for status, count := range statusCount {
		log.Printf("  %d: %d", status, count)
	}

	log.Println("\n请求主机 (Top 5):")
	topHosts := getTopEntries(hostCount, 5)
	for _, host := range topHosts {
		log.Printf("  %s: %d", host.key, host.count)
	}

	log.Println("==============================")
}

// entryCount 用于排序的辅助结构
type entryCount struct {
	key   string
	count int
}

// getTopEntries 获取计数最多的N个条目
func getTopEntries(m map[string]int, n int) []entryCount {
	entries := make([]entryCount, 0, len(m))
	for k, v := range m {
		entries = append(entries, entryCount{key: k, count: v})
	}

	// 简单排序
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].count > entries[i].count {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	if len(entries) > n {
		entries = entries[:n]
	}
	return entries
}