package core

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

/*
创建 Redis 布隆过滤器。

host - Redis 主机，如："localhost:6379"

auth - 密码

filterName - 过滤器名称
*/
func NewRedisBloom(host string, name string, auth string, filterName string) RedisBloomFilter {
	filter := new(RedisBloomFilter)
	filter.filterName = filterName
	var au *string = nil
	if auth != "" {
		au = &auth
	}
	filter.client = redisbloom.NewClient(host, name, au)
	return *filter
}

func (f RedisBloomFilter) Exists(key string) (exists bool, err error) {
	return f.client.Exists(f.filterName, key)
}

func (f RedisBloomFilter) Add(key string) (exists bool, err error) {
	return f.client.Add(f.filterName, key)
}

func (f RedisBloomFilter) Clear() (success bool, err error) {
	conn := f.client.Pool.Get()
	defer conn.Close()
	args := redis.Args{f.filterName}
	result, err := conn.Do("DEL", args...)
	return redis.Bool(result, err)
}
