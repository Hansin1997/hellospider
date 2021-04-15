package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Namespace     string              `json:"namespace"`
	Redis         RedisBloomConfig    `json:"redis"`
	RabbitMq      RabbitMqConfig      `json:"rabbitMq"`
	Elasticsearch ElasticsearchConfig `json:"elasticsearch"`
	Seeds         []string            `json:"seeds"`
	UserAgents    []string            `json:"userAgents"`
	Workers       int                 `json:"workers"`
	Accepts       []string            `json:"accepts"`
	// 优先级策略：0-9 表示优先级为常数，url-len 表示根据 URL 长度计算优先级（越短越优先），path-len 表示根据 URL 路径长度计算优先级（越短越优先）。
	Priority        string   `json:"priority"`
	Rules           Rule     `json:"rules"`
	ResponseHeaders []string `json:"responseHeaders"`
}

type RedisBloomConfig struct {
	Host string `json:"host"`
	Auth string `json:"auth"`
}

type ElasticsearchConfig struct {
	Address  []string `json:"address"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

type RabbitMqConfig struct {
	Url       string `json:"url"`
	Exchange  string `json:"exchange"`
	MaxLength int64  `json:"MaxLength"`
}

type Rule struct {
	Allows []string `json:"allows"`
	Forbid []string `json:"forbid"`
}

func loadConfig(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	return cfg, nil
}
