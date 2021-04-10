# Hello Spider 🕷
基于 ```Go``` 语言实现的分布式网页爬虫。

## 简介
此项目初衷在于学习 ```Go``` 语言，以及当前热门的 ```Elasticsearch``` 。

## 使用方法
### 命令行参数
```bash
$ go run . -h
Usage of hello-spider.exe:
  -config string
        配置文件路径。 (default "config.json")
  -reset
        开始前清空数据。
  -seed string
        种子 URL。
```
### 配置文件
修改配置文件 ```config.json``` 中各服务的地址、端口以及用户名及密码等。

```json
{
    "workers": 8,
    "seeds": [
        "https://baidu.com/"
    ],
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
    "accepts": [
        "text/html",
        "text/plain"
    ],
    "userAgents": [
        "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
        "Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
        "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
        "Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
        "Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
        "Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
        "Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
        "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
        "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
        "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52"
    ]
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

