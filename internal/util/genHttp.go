/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/23
  @desc: //TODO
**/

// Package util 生成Http响应内容
package util

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// HttpContext 生成http响应
func HttpContext(data string) []byte {
	l := len(data)
	ls := strconv.Itoa(l)
	httpData := "HTTP/1.1 200 OK\n" +

		//"Date: " + time.Now().String() + "\n" +
		"Content-Type: text/html; charset=UTF-8\n" +
		"Content-Length: " + ls + "\n\n" +
		"\r\n" +
		data
	//i := len(data)
	return []byte(httpData)
}

func ResponseToBytes(response *http.Response) ([]byte, error) {
	httpDataBuffer := ""
	httpDataBuffer += "HTTP/1.1 200 OK\n"
	for key, values := range response.Header {
		for _, value := range values {
			httpDataBuffer += fmt.Sprintf("%s: %s\r\n", key, value)
		}
	}
	httpDataBuffer += "\r\n"
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	res := append([]byte(httpDataBuffer), respBytes...)
	return res, err
}

func NewRequest(r *http.Request) (*http.Request, error) {
	req, _ := http.NewRequest(r.Method, r.RequestURI, r.Body)
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	//client := &http.Client{}
	//resp, err := client.Do(req)

	return req, nil
}

//
//func HttpResponse(data string) {
//	//response := &http.Response{
//	//	StatusCode:       http.StatusOK,
//	//	Proto:            "",
//	//	ProtoMajor:       0,
//	//	ProtoMinor:       0,
//	//	Header:           nil,
//	//	Body:             nil,
//	//	ContentLength:    0,
//	//	TransferEncoding: nil,
//	//	Close:            false,
//	//	Uncompressed:     false,
//	//	Trailer:          nil,
//	//	Request:          nil,
//	//	TLS:              nil,
//	//}
//	////err := response.Write()
//	//if err != nil {
//	//	return
//	//}
//
//}
