package core

import (
	"net/http"
	"regexp"
	"strings"
)

// SelectHeader 筛选 HTTP 头
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

// CheckUrl 检查 URL 是否有效
func CheckUrl(url string, allows []regexp.Regexp, forbid []regexp.Regexp) bool {
	if url == "" || strings.HasPrefix(url, "#") || strings.HasPrefix(url, "javascript") {
		return false
	}
	if len(allows) != 0 {
		match := false
		for _, allow := range allows {
			if allow.MatchString(url) {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	if len(forbid) != 0 {
		for _, f := range forbid {
			if f.MatchString(url) {
				return false
			}
		}
	}
	return true
}
