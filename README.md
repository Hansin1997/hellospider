# Hello Spider 🕷
基于 ```Go``` 语言实现的分布式网页爬虫。

## 简介
此项目初衷在于学习 ```Go``` 语言以及 ```Elasticsearch``` 。

## 使用方法
### 命令行参数
```bash
$ go run . -h
Usage of hello-spider:
  -config string
        File path of configuration. (default "config.json")
  -namespace string
        Namespace of task. (default "default")
  -reset
        Reset queue, storage and filter before begin task.
  -seed string
        The seeds URL is comma-separated. Such as: 'http://a.com/, http://b.com/'. And the seeds in the configuration file will be ignored.
```
### 配置文件
修改配置文件 ```config.json``` 中各服务的地址、端口以及用户名及密码等。

```json
{
    "namespace": "default",
    "workers": 8,
    "seeds": [
        "https://bing.com/"
    ],
    "redis": {
        "host": "localhost:6379",
        "auth": null
    },
    "rabbitMq": {
        "url": "amqp://guest:guest@localhost:5672/",
        "exchange": "spider"
    },
    "elasticsearch": {
        "address": [
            "http://localhost:9200"
        ],
        "username": "elastic",
        "password": "123456"
    },
    "accepts": [
        "text/html",
        "text/plain"
    ],
    "userAgents": [
        "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
        ...
    ]
}
```

## 相关技术
### 布隆过滤器
```RedisBloom``` 用于分布式下的 ```URL``` 去重。

[RedisBloom](https://github.com/RedisBloom/RedisBloom)

[redisbloom-go](https://github.com/RedisBloom/redisbloom-go)

### 消息队列
使用 ```RabbitMQ``` 存放待抓取的 ```URL```，并从队列获取 ```URL``` 进行消费。

[RabbitMQ Server](https://github.com/rabbitmq/rabbitmq-server)

[Go RabbitMQ Client Library](https://github.com/streadway/amqp)

### 数据存储
使用 ```Elasticsearch``` 存储网页数据。

[Elasticsearch](https://github.com/elastic/elasticsearch)

[go-elasticsearch](https://github.com/elastic/go-elasticsearch)

