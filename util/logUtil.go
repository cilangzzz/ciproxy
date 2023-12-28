// Package util 日记配置
package util

import (
	"log"
	"os"
)

// LogInit 日记初始化
// init log
// if path equal none print to terminal
// if path not none, try to create dir
func LogInit(path string) {
	if path == "" {
		println("log path not setting, print to terminal")
		return
	}
	_, err := os.Stat(path)
	if err != nil {
		println("log path not exist auto create")
		err := os.Mkdir("./log/", os.ModePerm)
		if err != nil {
			println("create log file")
			println(err)
			panic(err)
		}
	}
	filePath := "./" + path
	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		println("log file setting err")
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	//currentTime := time.Now().String()
	//log.SetPrefix("[" + currentTime + "]")
	log.SetPrefix("[CiProxy]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}
