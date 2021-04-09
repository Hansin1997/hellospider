package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	RedisHost       string  `json:"redisHost"`
	RedisAuth       *string `json:"redusAuth"`
	RedisClientName string  `json:"redisClientName"`
	FilterName      string  `json:"filerName"`
}

func loadConfig(path string) Config {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	cfg := Config{}
	json.Unmarshal(buf, cfg)
	return cfg
}

func initBloomFilter(config Config) BloomFilter {
	return newRedisBloom(config.RedisHost, config.RedisClientName, config.RedisAuth, config.FilterName)
}

func initFetcher(config Config) Fetcher {
	return newDefaultFetcher()
}

func reset(bloomFilter BloomFilter) {
	bloomFilter.Clear()
}

func main() {
	var _reset bool = false // 是否重置

	cfg := loadConfig("config.json") // 载入配置

	// var bloom BloomFilter = initBloomFilter(cfg) // 初始化布隆过滤器
	var fetcher Fetcher = initFetcher(cfg)

	if _reset {
		// reset(bloom) // 重置
	}

	doc, urls, err := fetcher.Fetch("https://gushi365.com/")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(doc)
	log.Println(urls)

}
