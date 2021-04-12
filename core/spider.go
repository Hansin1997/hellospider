package core

import (
	"log"
	"net/url"
	"strings"
	"sync"
)

// 爬虫
type Spider struct {
	Filter  BloomFilter
	Queue   Queue
	Fetcher Fetcher
	Storage Storage
}

// 重置任务
func (spider Spider) Reset() error {
	_, err := spider.Filter.Clear()
	if err != nil {
		return err
	}
	err = spider.Queue.Clear()
	if err != nil {
		return err
	}
	err = spider.Storage.Clear()
	if err != nil {
		return err
	}
	return nil
}

// 开始或者继续任务
func (spider Spider) Run(workers int) {
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := spider.Queue.Consume(func(content string) (bool, bool, error) {
				if !CheckUrl(content) {
					return true, false, nil
				}
				u, err := url.Parse(content)

				if err != nil {
					log.Printf("[Warning] Fail to fetch URL: %s \n%s\n", content, err)
					return true, false, nil
				}
				target := computeUrl(u)
				doc, urls, success, err := spider.Fetcher.Fetch(target) // 抓取网页
				if err != nil {
					log.Printf("[Warning] Fail to fetch URL: %s \n%s\n", content, err)
					return true, false, nil
				}
				if !success {
					return true, false, nil
				}
				err = spider.Storage.Save(doc) // 储存文档
				if err != nil {
					return false, true, err
				}
				log.Printf("[Info] %s [%s]\n", target, doc.Title)
				for _, newUrl := range urls {
					newUrl := strings.TrimSpace(newUrl)
					if !CheckUrl(newUrl) {
						continue
					}
					nu, err := url.Parse(newUrl)
					if err != nil {
						log.Printf("[Warning] Fail to push URL: %s \n%s\n", newUrl, err)
						continue
					}
					r := u.ResolveReference(nu)

					_, err = spider.Enqueue(r)

					if err != nil {
						return true, false, err
					}

				}
				return true, false, nil
			})
			if err != nil {
				log.Panicf("[Error] %s", err)
			}
		}()
	}
	wg.Wait()
}

// 将 URL 入队
func (spider Spider) Enqueue(target *url.URL) (success bool, err error) {
	newKey := computeKey(target)
	newUrl := computeUrl(target)
	if !CheckUrl(newUrl) {
		return false, nil
	}
	exists, err := spider.Filter.Exists(newKey) // 检查是否存在
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	err = spider.Queue.Publish(newUrl) // 入队
	if err != nil {
		return false, err
	}
	_, err = spider.Filter.Add(newKey) // 更新过滤器
	return true, err
}

// 简化 URL
func computeUrl(u *url.URL) string {
	if u.Host == "" {
		return ""
	}
	userStirng := strings.TrimSpace(u.User.String())
	if userStirng == "" {
		return u.Scheme + "://" + u.Host + u.RequestURI()
	} else {
		return u.Scheme + "://" + userStirng + "@" + u.Host + u.RequestURI()
	}
}

// 简化 URL （忽略协议）
func computeKey(u *url.URL) string {
	if u.Host == "" {
		return ""
	}
	userStirng := strings.TrimSpace(u.User.String())
	if userStirng == "" {
		return u.Host + u.RequestURI()
	} else {
		return userStirng + "@" + u.Host + u.RequestURI()
	}
}
