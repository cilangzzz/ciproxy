/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/6/12
  @desc: //TODO
**/

package serverProxy

import (
	"crypto/tls"
)

type (
	ServerProxy struct {
		Ip        string
		Port      string
		Method    string
		TLSConfig *tls.Config
	}
)
