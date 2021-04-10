# Hello Spider 🕷
基于 ```Go``` 语言实现的分布式网页爬虫。

## 简介
此项目初衷在于学习 ```Go``` 语言，以及当前热门的 ```Elasticsearch``` 。

## 开始
### 配置文件
修改配置文件 ```config.json``` 中各服务的地址、端口以及用户名及密码等。

```json
{
    "redis": {
        "host": "localhost:6379",
        "auth": null,
        "client": "spider",
        "filter": "spider"
    },
    "rabbitMq": {
        "url": "amqp://guest:guest@localhost:5672/",
        "exchange": "spider",
        "queue": "spider-work",
        "routingKey": "spider-work"
    },
    "elasticsearch": {
        "address": [
            "http://localhost:9200"
        ],
        "username": "elastic",
        "password": "123456",
        "index": "spider"
    },
    "seeds": [
        "https://qq.com"
    ],
    "reset": true
}
```

### 运行
```bash
go run .
```


## 相关技术
### 布隆过滤器
```RedisBloom``` 用于分布式下的 ```URL``` 去重。

[RedisBloom](https://github.com/RedisBloom/RedisBloom)

[redisbloom-go](https://github.com/RedisBloom/redisbloom-go)

### 消息队列
使用 ```RabbitMQ``` 存放待爬取的 ```URL```，并从队列获取 ```URL``` 并进行消费。

[RabbitMQ Server](https://github.com/rabbitmq/rabbitmq-server)

### 数据存储
使用 ```Elasticsearch``` 存储网页数据。

[Elasticsearch](https://github.com/elastic/elasticsearch)

