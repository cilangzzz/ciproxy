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

import "time"

// HttpContext 生成http响应
func HttpContext(data string) []byte {
	httpData := "HTTP/1.1 200 OK\n" +
		"Date: " + time.Now().String() + "\n" +
		"Content-Type: text/html; charset=UTF-8\n" +
		"Content-Length: 11\n\n" +
		data
	return []byte(httpData)
}
