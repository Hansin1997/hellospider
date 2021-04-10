package main

import (
	"log"
	"net/url"
	"strings"
)

func initBloomFilter(config Config) BloomFilter {
	return newRedisBloom(config.Redis.Host, config.Redis.Client, config.Redis.Auth, config.Redis.Filter)
}

func initFetcher(config Config) Fetcher {
	return newDefaultFetcher()
}

func initQueue(config Config) (Queue, error) {
	return newRbQueue(config.RabbitMq.Url, config.RabbitMq.Exchange, config.RabbitMq.Queue)
}

func reset(bloomFilter BloomFilter, queue Queue) {
	bloomFilter.Clear()
	queue.Clear()
}

func computeUrl(u *url.URL) string {
	userStirng := strings.TrimSpace(u.User.String())
	if userStirng == "" {
		return u.Scheme + "://" + u.Host + u.RequestURI()
	} else {
		return u.Scheme + "://" + userStirng + "@" + u.Host + u.RequestURI()
	}
}

func main() {
	var _reset bool = true // 是否重置

	cfg := loadConfig("config.json") // 载入配置

	filter := initBloomFilter(cfg) // 初始化布隆过滤器
	fetcher := initFetcher(cfg)
	queue, err := initQueue(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if _reset {
		reset(filter, queue) // 重置
	}

	queue.Consume(func(content string) (bool, error) {

		u, err := url.Parse(content)

		if err != nil {
			log.Println(err)
			return true, nil
		}
		target := computeUrl(u)
		doc, urls, success, err := fetcher.Fetch(target)
		if err != nil {
			log.Println(err)
			return true, nil
		}
		if !success {
			return true, nil
		}
		log.Printf("%s\t[%s]\n", target, doc.Title)
		for _, newUrl := range urls {
			newUrl = strings.TrimSpace(newUrl)
			if newUrl == "" || strings.HasPrefix(newUrl, "#") || strings.HasPrefix(newUrl, "javascript") {
				continue
			}
			nu, err := url.Parse(newUrl)
			if err != nil {
				log.Println(err)
				continue
			}
			r := u.ResolveReference(nu)
			newUrl = computeUrl(r)
			exists, err := filter.Exists(newUrl)
			if err != nil {
				return false, err
			}
			if exists {
				continue
			}
			err = queue.Publish(newUrl)
			if err != nil {
				return false, nil
			}
			filter.Add(newUrl)
		}
		return true, nil
	})
}
