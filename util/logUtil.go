package util

import (
	"log"
	"os"
)

func LogInit(path string) {
	_, err := os.Stat(path)
	if err != nil {
		println("log path not exist auto create")
		err := os.Mkdir("./log/", os.ModePerm)
		if err != nil {
			println(err)
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
