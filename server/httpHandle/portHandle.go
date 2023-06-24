/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/24
  @desc: //TODO
**/

package httpHandle

import "net/http"

type (
	PortHandle struct {
	}
)

func (s PortHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

}

func (s PortHandle) getPort(r *http.Request) {

}
