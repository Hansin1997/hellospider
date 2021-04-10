package main

import (
	"context"
	"log"
	"net/url"
	"strings"
)

func initBloomFilter(config Config) BloomFilter {
	log.Println("[Init BloomFilter]")
	return newRedisBloom(config.Redis.Host, config.Redis.Client, config.Redis.Auth, config.Redis.Filter)
}

func initFetcher(config Config) Fetcher {
	log.Println("[Init Fetcher]")
	return newDefaultFetcher()
}

func initQueue(config Config) (Queue, error) {
	log.Println("[Init Queue]")
	return newRbQueue(config.RabbitMq.Url, config.RabbitMq.Exchange, config.RabbitMq.Queue, config.RabbitMq.RoutingKey)
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
		seed = computeUrl(u)
		exists, err := filter.Exists(seed)
		if err != nil {
			return err
		}
		if !exists {
			err = queue.Publish(seed)
			if err != nil {
				return err
			}
			_, err = filter.Add(seed)
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

func main() {

	cfg := loadConfig("config.json") // 载入配置
	_reset := cfg.Reset              // 是否重置
	filter := initBloomFilter(cfg)   // 初始化布隆过滤器
	fetcher := initFetcher(cfg)
	queue, err := initQueue(cfg)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := initStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if _reset {
		reset(filter, queue, storage) // 重置
	}

	err = pushSeeds(filter, queue, cfg.Seeds)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("")
	log.Println("[Start]")
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
		err = storage.Save(doc)
		if err != nil {
			log.Println(err)
			return false, nil
		}
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

	log.Println("[Finished]")
}
