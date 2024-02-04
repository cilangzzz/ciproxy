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

// HttpContext 生成http响应
func HttpContext(data string) []byte {
	//l := len(data)
	//ls := strconv.Itoa(l)
	httpData := "HTTP/1.1 200 OK\n" +
		//"Date: " + time.Now().String() + "\n" +
		//"Content-Type: text/html; charset=UTF-8\n" +
		//"Content-Length: " + ls + "\n\n" +
		data
	//i := len(data)
	return []byte(httpData)
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
