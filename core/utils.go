package core

import (
	"net/http"
	"strings"
)

// 筛选 HTTP 头
func SelectHeader(header http.Header, allowHeaders map[string]bool) http.Header {
	result := http.Header{}
	for k, v := range allowHeaders {
		if !v {
			continue
		}
		val := header.Get(k)
		if val != "" {
			result.Set(k, val)
		}
	}
	return result
}

// 检查 URL 是否有效
func CheckUrl(url string) bool {
	if url == "" || strings.HasPrefix(url, "#") || strings.HasPrefix(url, "javascript") {
		return false
	}
	return true
}
