package core

import (
	"math"
	"net/http"
	"net/url"
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

// 计算消息的优先级，优先级从 0 - 9 递增
func GetPriority(content string) uint8 {
	u, err := url.Parse(content)
	var l int
	if err == nil {
		l = len(u.RequestURI())
	} else {
		l = len(content)
	}
	if l > 128 {
		return 0
	} else if l <= 1 {
		return 9
	} else {
		y := priorityFx(l)
		if y > 8 {
			y = 8
		} else if y < 0 {
			y = 0
		}
		return uint8(y)
	}
}

// 优先级函数 f(x)=(𝑒^((−(𝑥−340))/50))/100
func priorityFx(x int) int {
	fx := math.Pow(math.E, -(float64(x)-340)/50.0) / 100.0
	return int(math.Floor(0.5 + fx)) // 四舍五入
}
