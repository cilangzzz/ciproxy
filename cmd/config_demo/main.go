/**
 * @file config_demo/main.go
 * @brief 配置文件使用示例
 * @description 展示如何使用JSON配置文件配置代理服务器
 *
 * 功能特点:
 *   - 支持JSON格式配置文件
 *   - 支持配置验证
 *   - 支持命令行参数覆盖配置文件
 *   - 支持生成默认配置文件
 *
 * 使用方法:
 *   # 使用默认配置
 *   go run main.go
 *
 *   # 使用指定配置文件
 *   go run main.go -config config.json
 *
 *   # 生成默认配置文件
 *   go run main.go -generate-config default.json
 *
 *   # 命令行参数覆盖配置
 *   go run main.go -config config.json -ip 0.0.0.0 -port 9090
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
	configPath := flag.String("config", "", "配置文件路径")
	generateConfig := flag.String("generate-config", "", "生成默认配置文件并退出")
	ip := flag.String("ip", "", "覆盖配置中的IP地址")
	port := flag.String("port", "", "覆盖配置中的端口")
	flag.Parse()

	// 生成默认配置文件
	if *generateConfig != "" {
		config := ciproxy.DefaultConfig
		if err := config.Save(*generateConfig); err != nil {
			log.Fatalf("保存配置失败: %v", err)
		}
		log.Printf("默认配置已保存到: %s", *generateConfig)
		return
	}

	// 加载配置
	var config *ciproxy.ProxyConfig
	var err error

	if *configPath != "" {
		config, err = ciproxy.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}
		log.Printf("已加载配置文件: %s", *configPath)
	} else {
		config = &ciproxy.DefaultConfig
		log.Println("使用默认配置")
	}

	// 命令行参数覆盖配置文件
	if *ip != "" {
		config.IP = *ip
	}
	if *port != "" {
		config.Port = *port
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 打印配置信息
	printConfig(config)

	// 使用配置创建服务器
	server := ciproxy.New(ciproxy.WithConfig(config))

	// 启动服务器
	go func() {
		log.Printf("代理服务器启动: %s:%s", config.IP, config.Port)
		log.Printf("代理模式: %s", config.Method)
		if config.Features.EnableTrafficCapture {
			log.Println("流量捕获: 已启用")
		}
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

// printConfig 打印配置信息
func printConfig(config *ciproxy.ProxyConfig) {
	log.Println("========== 当前配置 ==========")
	log.Printf("监听地址: %s:%s", config.IP, config.Port)
	log.Printf("代理模式: %s", config.Method)
	log.Printf("协议: %s", config.Protocol)
	if config.LogPath != "" {
		log.Printf("日志路径: %s", config.LogPath)
	}
	log.Printf("TLS证书: %s", config.TLS.CertPath)
	log.Printf("连接超时: %v", config.Timeout.Connect)
	log.Printf("读取超时: %v", config.Timeout.Read)
	log.Printf("写入超时: %v", config.Timeout.Write)
	log.Println("功能开关:")
	log.Printf("  流量捕获: %v", config.Features.EnableTrafficCapture)
	log.Printf("  流量重放: %v", config.Features.EnableTrafficReplay)
	log.Printf("  隧道加密: %v", config.Features.EnableTunnelCrypt)
	log.Println("==============================")
}