package util

import (
	"log"
	"os"
)

func LogInit(path string) {
	_, err := os.Stat(path)
	println()
	if err != nil {
		println("log file setting err")
		err := os.Mkdir("./log/", os.ModePerm)
		if err != nil {
			println(err)
		}
	}
	filePath := "./" + path
	println(filePath)
	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		println("log file setting err")
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[qSkipTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}
