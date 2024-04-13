/**
  @creator: cilang
  @qq: 1019383856
  @bili: https://space.bilibili.com/433915419
  @gitee: https://gitee.com/OpencvLZG
  @github: https://github.com/OpencvLZG
  @since: 2023/12/25
  @desc: //TODO
**/

package main

import (
	"flag"
	"github.com/cilangzzz/ciproxy/internal/util"
)

func main() {
	fileType := flag.String("fileType", "crt", "file format to save")
	organization := flag.String("organization", "www.cilang.buzz", "certificate organization")
	country := flag.String("country", "cn", "certificate country")
	province := flag.String("province", "www.cilang.buzz", "certificate province")
	locality := flag.String("locality", "GuangZhou", "locality")
	organizationalUnit := flag.String("organizationalUnit", "software", "certificate organizationalUnit")
	commonName := flag.String("commonName", "localhost", "certificate commonName")
	dnsDomain := flag.String("dnsDomain", "docker.cilang.buzz", "certificate dnsDomain")
	flag.Parse()
	util.GenerateCert(*fileType, *organization, *country, *province, *locality, *organizationalUnit, *commonName, *dnsDomain)
}
