package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Redis         RedisBloomConfig    `json:"redis"`
	RabbitMq      RabbitMqConfig      `json:"rabbitMq"`
	Elasticsearch ElasticsearchConfig `json:"elasticsearch"`
	Seeds         []string            `json:"seeds"`
	Reset         bool                `json:"reset"`
	Workers       int                 `json:"workers"`
	Accepts       []string            `json:"accepts"`
}

type RedisBloomConfig struct {
	Host   string `json:"host"`
	Auth   string `json:"auth"`
	Client string `json:"client"`
	Filter string `json:"filter"`
}

type ElasticsearchConfig struct {
	Address  []string `json:"address"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Index    string   `json:"index"`
}

type RabbitMqConfig struct {
	Url        string `json:"url"`
	Exchange   string `json:"exchange"`
	Queue      string `json:"queue"`
	RoutingKey string `json:"routingKey"`
}

func loadConfig(path string) Config {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	cfg := Config{}
	json.Unmarshal(buf, &cfg)
	return cfg
}
