/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/6/11
  @desc: //TODO
**/

// Package trafficHandle 流量转发
package trafficHandle

import (
	"io"
	"log"
)

func errLog(msg string, err error) {
	log.Println("proxyHandle:" + msg + " err:" + err.Error())
}

// cryptTraffic 流量加密
func cryptTraffic() {

}

// decryptTraffic 流量解密
func decryptTraffic() {

}

// TunnelTransfer 加密流量转发
func TunnelTransfer() {

}

// Transfer traffic transfer 流量Io转发
func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func(destination io.WriteCloser) {
		err := destination.Close()
		if err != nil {
			//errLog("close io writer failed", err)
		}
	}(destination)
	defer func(source io.ReadCloser) {
		err := source.Close()
		if err != nil {
			//errLog("close io writer failed", err)
		}
	}(source)
	_, err := io.Copy(destination, source)
	if err != nil {
		errLog("copy data transfer failed", err)
	}
}
