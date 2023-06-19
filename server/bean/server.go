/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/19
  @desc: //TODO
**/

package bean

type Server struct {
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Method   string `json:"method,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}
