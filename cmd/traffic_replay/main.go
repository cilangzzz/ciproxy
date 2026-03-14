/**
 * @file traffic_replay/main.go
 * @brief 流量重放工具示例
 * @description 展示如何使用CiProxy重放捕获的HTTP请求
 *
 * 功能特点:
 *   - 支持从JSON文件导入流量
 *   - 支持并行重放请求
 *   - 支持请求修改规则
 *   - 支持重放前后回调
 *
 * 使用方法:
 *   # 重放单个请求
 *   go run main.go -file traffic.json -id <request-id>
 *
 *   # 重放所有请求
 *   go run main.go -file traffic.json -all
 *
 *   # 并行重放（10并发）
 *   go run main.go -file traffic.json -all -concurrency 10
 *
 * @author cilang
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/opencvlzg/ciproxy/pkg/transfer"
)

func main() {
	// 命令行参数解析
	file := flag.String("file", "", "流量文件路径 (JSON格式)")
	id := flag.String("id", "", "要重放的请求ID")
	all := flag.Bool("all", false, "重放所有请求")
	concurrency := flag.Int("concurrency", 1, "并行重放并发数")
	timeout := flag.Duration("timeout", 30*time.Second, "请求超时时间")
	verbose := flag.Bool("v", false, "详细输出")
	flag.Parse()

	if *file == "" {
		log.Fatal("请指定流量文件路径: -file <path>")
	}

	// 加载流量数据
	store := transfer.NewTrafficStore()
	if err := loadTrafficFile(store, *file); err != nil {
		log.Fatalf("加载流量文件失败: %v", err)
	}
	log.Printf("已加载 %d 条流量记录", store.Len())

	if store.Len() == 0 {
		log.Fatal("流量文件中没有数据")
	}

	// 创建重放器
	replayer := transfer.NewReplayer(
		transfer.WithTimeout(*timeout),
		transfer.WithInsecureSkipVerify(true),
	)

	// 设置详细输出回调
	if *verbose {
		replayer.SetBeforeReplay(func(req *http.Request) error {
			log.Printf("[REPLAY] %s %s", req.Method, req.URL)
			return nil
		})
	}

	// 执行重放
	if *id != "" {
		// 重放单个请求
		replaySingle(replayer, store, *id)
	} else if *all {
		// 重放所有请求
		if *concurrency > 1 {
			replayAllParallel(replayer, store, *concurrency)
		} else {
			replayAll(replayer, store)
		}
	} else {
		// 列出可用的请求
		listRequests(store)
		log.Fatal("请指定重放模式: -id <id> 或 -all")
	}
}

// loadTrafficFile 从文件加载流量数据
func loadTrafficFile(store *transfer.TrafficStore, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return store.Import(data)
}

// replaySingle 重放单个请求
func replaySingle(replayer *transfer.Replayer, store *transfer.TrafficStore, id string) {
	entry, ok := store.Get(id)
	if !ok {
		log.Fatalf("找不到ID为 %s 的请求", id)
	}

	log.Printf("重放请求: %s %s", entry.Request.Method, entry.Request.URL)
	startTime := time.Now()

	resp, err := replayer.Replay(entry)
	if err != nil {
		log.Fatalf("重放失败: %v", err)
	}

	log.Printf("重放成功: ID=%s, Status=%d, 耗时=%v",
		id, resp.Status, time.Since(startTime))
}

// replayAll 重放所有请求
func replayAll(replayer *transfer.Replayer, store *transfer.TrafficStore) {
	log.Printf("开始重放所有请求...")
	startTime := time.Now()

	responses, errs := replayer.ReplayAll(store)

	log.Printf("重放完成: 成功 %d, 失败 %d, 耗时 %v",
		len(responses), len(errs), time.Since(startTime))

	// 打印错误
	for i, err := range errs {
		log.Printf("错误 %d: %v", i+1, err)
	}

	// 打印状态码分布
	printStatusDistribution(responses)
}

// replayAllParallel 并行重放所有请求
func replayAllParallel(replayer *transfer.Replayer, store *transfer.TrafficStore, concurrency int) {
	log.Printf("开始并行重放 (并发: %d)...", concurrency)
	startTime := time.Now()

	responses, errs := replayer.ReplayParallel(store, concurrency)

	log.Printf("并行重放完成: 成功 %d, 失败 %d, 并发 %d, 耗时 %v",
		len(responses), len(errs), concurrency, time.Since(startTime))

	// 打印状态码分布
	printStatusDistribution(responses)
}

// listRequests 列出所有请求
func listRequests(store *transfer.TrafficStore) {
	entries := store.List()
	fmt.Println("\n可用的请求:")
	fmt.Println("==========================================")
	for i, entry := range entries {
		if i >= 20 {
			fmt.Printf("... 还有 %d 条记录\n", len(entries)-20)
			break
		}
		fmt.Printf("ID: %s | %s %s\n",
			entry.ID,
			entry.Request.Method,
			entry.Request.URL)
	}
	fmt.Println("==========================================")
}

// printStatusDistribution 打印状态码分布
func printStatusDistribution(responses []*transfer.CapturedResponse) {
	if len(responses) == 0 {
		return
	}

	statusCount := make(map[int]int)
	for _, resp := range responses {
		statusCount[resp.Status]++
	}

	fmt.Println("\n状态码分布:")
	for status, count := range statusCount {
		fmt.Printf("  %d: %d\n", status, count)
	}
}