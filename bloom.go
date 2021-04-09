package main

import (
	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/gomodule/redigo/redis"
)

// 布隆过滤器接口
type BloomFilter interface {
	// 判断 Key 是否存在
	Exists(key string) (exists bool, err error)
	// 添加 Key
	Add(key string) (exists bool, err error)
	// 清空
	Clear() (success bool, err error)
}

// Redis 布隆过滤器
type RedisBloomFilter struct {
	client     *redisbloom.Client
	filterName string
}

// 创建 Redis 布隆过滤器。
//
// host - Redis 主机，如："localhost:6379"
//
// name - Redis Client Name
//
// auth - 密码
//
// filterName - 过滤器名称
func newRedisBloom(host string, name string, auth *string, filterName string) RedisBloomFilter {
	filter := new(RedisBloomFilter)
	filter.filterName = filterName
	filter.client = redisbloom.NewClient(host, name, auth)
	return *filter
}

func (instance RedisBloomFilter) Exists(key string) (exists bool, err error) {
	return instance.client.Exists(instance.filterName, key)
}

func (instance RedisBloomFilter) Add(key string) (exists bool, err error) {
	return instance.client.Add(instance.filterName, key)
}

func (instance RedisBloomFilter) Clear() (success bool, err error) {
	conn := instance.client.Pool.Get()
	defer conn.Close()
	args := redis.Args{instance.filterName}
	result, err := conn.Do("DEL", args...)
	return redis.Bool(result, err)
}
