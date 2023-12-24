/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/12
  @desc: //TODO
**/

package trafficForward

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

func HandleServerConnect(client net.Conn, proxyHost string, tlsConfig *tls.Config) {
	//buf := make([]byte, 1024)
	//_, err := client.Read(buf)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	target, err := tls.Dial("tcp", proxyHost, tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	go transfer(target, client)
	go transfer(client, target)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
