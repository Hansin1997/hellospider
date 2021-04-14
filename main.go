package main

import (
	"context"
	"flag"
	"hellospider/core"
	"log"
	"net/url"
	"strings"
	"time"
)

// 初始化布隆过滤器
func initBloomFilter(config Config) core.BloomFilter {
	return core.NewRedisBloom(config.Redis.Host, "spider-"+config.Namespace, config.Redis.Auth, "filter:"+config.Namespace)
}

// 初始化抓取器
func initFetcher(config Config) core.Fetcher {
	return core.NewDefaultFetcher(config.Accepts, config.UserAgents, config.ResponseHeaders)
}

// 初始化消息队列
func initQueue(config Config) (core.Queue, error) {
	return core.NewRbQueue(config.RabbitMq.Url,
		config.RabbitMq.Exchange,
		"spider-"+config.Namespace,
		"spider/"+config.Namespace,
		config.Workers*8,
		config.RabbitMq.MaxLength,
		config.Priority)
}

// 初始化存储器
func initStorage(config Config) (core.Storage, error) {
	return core.NewElasticsearchStorage(config.Elasticsearch.Address,
		config.Elasticsearch.Username,
		config.Elasticsearch.Password,
		"spider-"+config.Namespace,
		context.Background())
}

// 初始化全部组件
func initAll(config Config) (spider *core.Spider, error error) {
	filter := initBloomFilter(config)
	fetcher := initFetcher(config)
	queue, err := initQueue(config)
	if err != nil {
		return nil, err
	}
	storage, err := initStorage(config)
	if err != nil {
		return nil, err
	}
	log.Println("[Info] Finish initialize.")
	return &core.Spider{
		Filter:  filter,
		Queue:   queue,
		Storage: storage,
		Fetcher: fetcher,
	}, nil
}

func main() {

	configFile := flag.String("config", "config.json", "File path of configuration.")
	seed := flag.String("seed", "", "The seeds URL is comma-separated. Such as: 'https://a.com/, https://b.com/'. And the seeds in the configuration file will be ignored.")
	reset := flag.Bool("reset", false, "Reset queue, storage and filter before begin task.")
	namespace := flag.String("namespace", "", "Namespace of task.")
	priority := flag.String("priority", "", "Priority policy: 0-9 means that the priority is constant, url-len means that the priority is calculated according to the length of the URL (the shorter the priority), path-len means that the priority is calculated according to the length of the URL path (the shorter the priority).")

	flag.Parse()
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("[Error] Fail to load config!\n%s\n", err)
	}

	ns := strings.TrimSpace(*namespace)
	if config.Namespace == "" || ns != "" {
		config.Namespace = ns
	}
	log.Printf("[Info] Using namespace [%s]\n", config.Namespace)

	pty := strings.TrimSpace(*priority)
	if config.Priority == "" || pty != "" {
		config.Priority = pty
	}
	log.Printf("[Info] Using priority policy [%s]\n", config.Priority)

	if *reset {
		for i := 5; i >= 0; i-- {
			log.Printf("[Warining] ⚠ Reset flag is true, clear namespace [%s] in %d seconds. ⚠ ", config.Namespace, i)
			time.Sleep(time.Second)
		}
	}

	spider, err := initAll(*config)
	if err != nil {
		log.Fatalf("[Error] Fail to initialize!\n%s\n", err)
	}
	if *reset {
		err = spider.Reset()
		if err != nil {
			log.Fatalf("[Error] Fail to reset!\n%s", err)
		}
		log.Println("[Info] Finish reset.")
	}

	if *seed != "" {
		seeds := strings.Split(*seed, ",")
		config.Seeds = seeds
	}
	for _, s := range config.Seeds {
		u, err := url.Parse(s)
		if err != nil {
			log.Printf("[Warning] Fail to push seed: %s\n%s\n", s, err)
			continue
		}
		success, err := spider.Enqueue(u)
		if err != nil {
			log.Printf("[Warning] Fail to push seed: %s\n%s\n", s, err)
		} else if success {
			log.Printf("[Info] Push seed: %s \n", s)
		} else {
			log.Printf("[Warning] Fail to push seed: %s\n", s)
		}
	}
	spider.Run(config.Workers)
}
