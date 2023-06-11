package util

import (
	"net"
	"regexp"
)

func GetHttpsHostRegex(url string) string {
	portReg := regexp.MustCompile(`(:)\d{2,5}`)
	port := portReg.FindString(url)
	reg := regexp.MustCompile(`(?i)(http://|https://|\/|:([0-9]+))`)
	url = reg.ReplaceAllString(url, "")
	switch port {
	case "443":
		url += ":443"
	case "":
		url += ":443"
	default:
		url += port
	}

	return url
}

func GetHttpHostRegex(url string) string {
	reg := regexp.MustCompile(`(?i)(http://|https://|\/|:80)`)
	url = reg.ReplaceAllString(url, "")
	url += ":80"
	return url
}

func IsValidHost(host string) bool {
	// Parse as IP address
	if ip := net.ParseIP(host); ip != nil {
		return true
	}

	// Parse as hostname
	if _, err := net.LookupHost(host); err == nil {
		return true
	}

	// Parse as URL
	//if _, err := url.Parse(host); err == nil {
	//	return true
	//}

	return false
}
