/**
  @author: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @since: 2023/5/21
  @desc: //TODO
**/

package main

import "trafficForward/server/serve"

func main() {
	proxyServe := serve.ProxyServe{}
	proxyServe.Start()
}
