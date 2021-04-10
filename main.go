package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"strings"
	"sync"
)

func initBloomFilter(config Config) BloomFilter {
	log.Println("[Init BloomFilter]")
	return newRedisBloom(config.Redis.Host, config.Redis.Client, config.Redis.Auth, config.Redis.Filter)
}

func initFetcher(config Config) Fetcher {
	log.Println("[Init Fetcher]")
	return newDefaultFetcher(config.Accepts, config.UserAgents)
}

func initQueue(config Config) (Queue, error) {
	log.Println("[Init Queue]")
	return newRbQueue(config.RabbitMq.Url,
		config.RabbitMq.Exchange,
		config.RabbitMq.Queue,
		config.RabbitMq.RoutingKey,
		config.Workers*8)
}

func initStorage(config Config) (Storage, error) {
	log.Println("[Init Storage]")
	return newElasticsearchStorage(config.Elasticsearch.Address,
		config.Elasticsearch.Username,
		config.Elasticsearch.Password,
		config.Elasticsearch.Index,
		context.Background())
}

func reset(bloomFilter BloomFilter, queue Queue, storage Storage) {
	_, err := bloomFilter.Clear()
	if err != nil {
		log.Printf("[Reset BloomFilter Error] %s\n", err.Error())
	} else {
		log.Println("[Reset BloomFilter Success]")
	}
	err = queue.Clear()
	if err != nil {
		log.Printf("[Reset Queue Error] %s\n", err.Error())
	} else {
		log.Println("[Reset Queue Success]")
	}
	err = storage.Clear()
	if err != nil {
		log.Printf("[Reset Storage Error] %s\n", err.Error())
	} else {
		log.Println("[Reset Storage Success]")
	}
}

func computeUrl(u *url.URL) string {
	userStirng := strings.TrimSpace(u.User.String())
	if userStirng == "" {
		return u.Scheme + "://" + u.Host + u.RequestURI()
	} else {
		return u.Scheme + "://" + userStirng + "@" + u.Host + u.RequestURI()
	}
}

func computeKey(u *url.URL) string {
	userStirng := strings.TrimSpace(u.User.String())
	if userStirng == "" {
		return u.Host + u.RequestURI()
	} else {
		return userStirng + "@" + u.Host + u.RequestURI()
	}
}

func pushSeeds(filter BloomFilter, queue Queue, seeds []string) error {
	if len(seeds) == 0 {
		return nil
	}
	for _, seed := range seeds {
		u, err := url.Parse(seed)
		if err != nil {
			log.Println(err)
			continue
		}
		key := computeKey(u)
		seed = u.Scheme + "://" + key
		exists, err := filter.Exists(key)
		if err != nil {
			return err
		}
		if !exists {
			err = queue.Publish(seed)
			if err != nil {
				return err
			}
			_, err = filter.Add(key)
			if err != nil {
				log.Println(err)
			}
			log.Printf("[Push Seed] %s\n", seed)
		} else {
			log.Printf("[Pass Seed] %s\n", seed)
		}
	}
	return nil
}

func handle(content string, filter BloomFilter, fetcher Fetcher, storage Storage, queue Queue) (bool, error) {
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
	err = storage.Save(doc)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	log.Printf("[Save] %s [%s]\n", target, doc.Title)
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
		newKey := computeKey(r)
		newUrl = r.Scheme + "://" + newKey
		exists, err := filter.Exists(newKey) // 检查是否存在
		if err != nil {
			return false, err
		}
		if exists {
			continue
		}
		err = queue.Publish(newUrl) // 入队
		if err != nil {
			return false, nil
		}
		filter.Add(newKey) // 更新过滤器
	}
	return true, nil
}

func main() {

	_reset := flag.Bool("reset", false, "开始前清空数据。")                // 是否重置
	_configFile := flag.String("config", "config.json", "配置文件路径。") // 是否重置
	_seed := flag.String("seed", "", "种子 URL。")
	flag.Parse()
	cfg := loadConfig(*_configFile) // 载入配置
	if cfg.Workers < 1 {
		cfg.Workers = 1
	}
	filter := initBloomFilter(cfg) // 初始化布隆过滤器
	fetcher := initFetcher(cfg)
	queue, err := initQueue(cfg)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := initStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if *_reset {
		reset(filter, queue, storage) // 重置
	}

	if strings.TrimSpace(*_seed) != "" {
		err = pushSeeds(filter, queue, []string{*_seed})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = pushSeeds(filter, queue, cfg.Seeds)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("")
	log.Println("[Start]")

	wg := sync.WaitGroup{}
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func(queue Queue) {
			err := queue.Consume(func(content string) (bool, error) {
				return handle(content, filter, fetcher, storage, queue)
			})
			if err != nil {
				log.Fatal(err)
			} else {
				wg.Done()
			}
		}(queue)
	}
	wg.Wait()
	log.Println("[Finished]")
}
