package util

import "regexp"

func GetHttpsHostRegex(url string) string {
	reg := regexp.MustCompile(`(?i)(http://|https://|\/|:([0-9]+))`)
	url = reg.ReplaceAllString(url, "")
	url += ":443"
	return url
}

func GetHttpHostRegex(url string) string {
	reg := regexp.MustCompile(`(?i)(http://|https://|\/|:80)`)
	url = reg.ReplaceAllString(url, "")
	url += ":80"
	return url
}
