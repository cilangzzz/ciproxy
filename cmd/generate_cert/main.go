/**
 * @file generate_cert/main.go
 * @brief TLS证书生成工具
 * @description 用于生成MITM代理所需的CA证书
 *
 * 功能特点:
 *   - 生成自签名CA证书
 *   - 支持自定义证书信息
 *   - 支持多种输出格式
 *
 * 使用方法:
 *   # 生成默认证书
 *   go run main.go
 *
 *   # 生成指定域名的证书
 *   go run main.go -commonName localhost -dnsDomain "example.com,*.example.com"
 *
 *   # 生成PEM格式证书
 *   go run main.go -fileType pem -organization "My Organization"
 *
 * 生成的证书文件:
 *   - root.crt  根证书（需要安装到系统信任列表）
 *   - private.pem  私钥文件（保密）
 *
 * 安装证书:
 *   Windows: 双击 root.crt -> 安装证书 -> 本地计算机 -> 受信任的根证书颁发机构
 *   macOS:   sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain root.crt
 *   Linux:   sudo cp root.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates
 *
 * 注意事项:
 *   1. 生成的CA证书需要安装到操作系统或浏览器的信任列表
 *   2. 私钥文件请妥善保管，不要泄露
 *   3. 生产环境建议使用正规CA机构签发的证书
 *
 * @author cilang
 */

package main

import (
	"flag"
	"log"

	"github.com/opencvlzg/ciproxy/pkg/util"
)

func main() {
	// 证书参数
	fileType := flag.String("fileType", "crt", "输出文件格式 (crt/pem)")
	organization := flag.String("organization", "www.cilang.buzz", "组织名称")
	country := flag.String("country", "cn", "国家代码 (如: cn, us, uk)")
	province := flag.String("province", "GuangDong", "省份")
	locality := flag.String("locality", "GuangZhou", "城市")
	organizationalUnit := flag.String("organizationalUnit", "software", "组织单位")
	commonName := flag.String("commonName", "localhost", "通用名称 (CN)")
	dnsDomain := flag.String("dnsDomain", "docker.cilang.buzz", "DNS域名 (多个用逗号分隔)")

	flag.Parse()

	log.Println("========== 证书生成工具 ==========")
	log.Printf("组织 (O): %s", *organization)
	log.Printf("国家 (C): %s", *country)
	log.Printf("省份 (ST): %s", *province)
	log.Printf("城市 (L): %s", *locality)
	log.Printf("组织单位 (OU): %s", *organizationalUnit)
	log.Printf("通用名称 (CN): %s", *commonName)
	log.Printf("DNS域名: %s", *dnsDomain)
	log.Println("==================================")

	// 生成证书
	util.GenerateCert(
		*fileType,
		*organization,
		*country,
		*province,
		*locality,
		*organizationalUnit,
		*commonName,
		*dnsDomain,
	)

	log.Println("证书生成完成!")
	log.Println("请将 root.crt 安装到系统/浏览器的信任列表")
	log.Println("")
	log.Println("安装方法:")
	log.Println("  Windows: 双击 root.crt -> 安装证书 -> 本地计算机 -> 受信任的根证书颁发机构")
	log.Println("  macOS:   sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain root.crt")
	log.Println("  Linux:   sudo cp root.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates")
}