package core

import (
	"math"
	"net/url"
	"strconv"
)

func GetPriorityFunc(policy string) (func(input string) uint8, error) {
	switch policy {
	case "0":
		fallthrough
	case "1":
		fallthrough
	case "2":
		fallthrough
	case "3":
		fallthrough
	case "4":
		fallthrough
	case "5":
		fallthrough
	case "6":
		fallthrough
	case "7":
		fallthrough
	case "8":
		fallthrough
	case "9":
		val, err := strconv.Atoi(policy)
		if err != nil {
			return nil, err
		}
		return func(input string) uint8 {
			return uint8(val)
		}, nil
		break
	case "url-len":
		return getPriorityByUrlLength, nil
	case "path-len":
		return getPriorityByUrlPathLength, nil
	}
	return nil, nil
}

// 根据 URL 长度计算消息的优先级，优先级从 0 - 9 递增
func getPriorityByUrlLength(content string) uint8 {
	l := len(content)
	if l > 128 {
		return 0
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

// 根据 URL 路径长度计算消息的优先级，优先级从 0 - 9 递增
func getPriorityByUrlPathLength(content string) uint8 {
	u, err := url.Parse(content)
	var l int
	if err == nil {
		l = len(u.RequestURI())
	} else {
		l = len(content)
	}
	if l > 128 {
		return 0
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
