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
	Url      string `json:"url"`
	Exchange string `json:"exchange"`
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
