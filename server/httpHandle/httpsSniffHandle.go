package httpHandle

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"trafficForward/server/middle"
	"trafficForward/server/util"
)

func HandleClientConnect(client net.Conn) {
	buf := make([]byte, 1024)
	_, err := client.Read(buf)

	host, err := util.ParseUrl(buf)
	println(string(buf))
	log.Println(host)
	if err != nil {
		return
	}
	var method, url string
	_, err = fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &url)
	if err != nil {
		return
	}

	if false {
		log.Println(buf)
		return
	}
	httpsSniffHandle(host, client)

}

func httpsSniffHandle(host string, client net.Conn) {
	tlsConfig := util.TLSUtil{Organization: "CiproxyOrganization"}
	cert, err := tlsConfig.GenCertificate()
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	target, err := tls.Dial("tcp", host, config)
	if err != nil {
		println(host)
		log.Println(err)
		return
	}
	_, err = client.Write([]byte("HTTP/1.1 200 Connection Established \r\n\r\n"))
	if err != nil {
		log.Println(err)
		return
	}
	mdManage := middle.MdManage
	for _, handle := range mdManage.HandleChain {
		handle(client, target)
	}

	go httpsSniffTransforward(client, target)
	go httpsSniffTransforward(target, client)
}

func httpsSniffTransforward(source net.Conn, destination net.Conn) {
	defer func(source net.Conn) {
		err := source.Close()
		if err != nil {
			log.Println(err)
		}
	}(source)
	defer func(destination net.Conn) {
		err := destination.Close()
		if err != nil {
			log.Println(err)
		}
	}(destination)
	buf := make([]byte, 32*1024)
	for {
		n, err := source.Read(buf)
		if err != nil {
			log.Println(err)
		}
		for n == len(buf) {
			buf = append(buf, make([]byte, 1024)...)
			n, err = source.Read(buf[len(buf)-1024:])
			if err != nil {
				log.Println(err)
			}
		}

		_, err = destination.Write(buf[:n])
		if err != nil {
			log.Println(err)
		}
	}
}
