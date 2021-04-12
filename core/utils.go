package core

import "strings"

func CheckUrl(url string) bool {
	if url == "" || strings.HasPrefix(url, "#") || strings.HasPrefix(url, "javascript") {
		return false
	}
	return true
}
