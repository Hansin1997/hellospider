# Hello Spider 🕷

基于 ```Go``` 语言的分布式网页爬虫。

*爬取数据后可使用 
[Hello Search](https://github.com/Hansin1997/hellosearch)
进行搜索。*

## 简介

> 此项目初衷在于学习 ```Go``` 语言以及 ```Elasticsearch``` 。

### 基本原理

1. 将种子 URL 入队。
2. 从队列获取待抓取 URL。
3. 判断此 URL 是否有效。（包括检查 URL 是否已经被抓取过）
4. 抓取网页信息，并记录此 URL 已经被抓取。
5. 储存网页摘要，并将此网页的超链接（URL）进行入队。
6. 执行步骤 ``` 2```。

### 实现概要

* 使用 ```RabbitMQ``` 消息队列实现分布式下的 URL 队列的持久化、优先级等。
* 使用 ```RedisBloom``` 的布隆过滤器实现分布式下的 URL 去重。
* 使用 ```goquery``` 解析 HTML。
* 使用 ```Elasticsearch``` 储存数据。

另外可以使用 ```Kibana``` 对抓取状况进行实时的可视化分析。

### Kibana 可视化

![Kibana可视化](docs/img/kibana.png?raw=true)

## 相关环境

### 构建环境

* ```go version go1.16.3```

### 运行环境

* ```RabbitMQ 3.8.12```
* ```RedisBloom latest```

## 使用方法

### 命令行参数

```bash
$ go run . -h
Usage of hellospider:
  -config string
        File path of configuration. (default "config.json")
  -namespace string
        Namespace of task.
  -priority string
        Priority policy: 0-9 means that the priority is constant, url-len means that the priority is calculated according to the length of the URL (the shorter the priority),
path-len means that the priority is calculated according to the length of the URL path (the shorter the priority).
  -reset
        Reset queue, storage and filter before begin task.
  -seed string
        The seeds URL is comma-separated. Such as: 'https://a.com/, https://b.com/'. And the seeds in the configuration file will be ignored.
```

* ```-config``` 配置文件路径
* ```-namespace``` 命名空间（区分不同任务）
* ```-priority``` 优先级策略：0-9 表示优先级为常数，url-len 表示根据 URL 长度计算优先级（越短越优先），path-len 表示根据 URL 路径长度计算优先级（越短越优先）。
* ```-reset``` 开始前是否重置命名空间
* ```-seed``` 替换配置文件中的 URL 种子（英文逗号分隔）

### 配置文件

修改配置文件 ```config.json``` 中各服务的地址、端口以及用户名及密码等。

```json
{
  "namespace": "default",
  "workers": 8,
  "priority": "path-len",
  "seeds": [
    "https://bing.com/"
  ],
  "redis": {
    "host": "localhost:6379",
    "auth": null
  },
  "rabbitMq": {
    "url": "amqp://guest:guest@localhost:5672/",
    "exchange": "spider",
    "maxLength": 999999999
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
    "..."
  ],
  "responseHeaders": [
    "Content-Type",
    "Content-Length",
    "Content-Language",
    "Server",
    "X-Powered-By"
  ]
}
```

### 运行

#### 源代码

```bash
go run .
```

### 二进制可执行文件

```bash
./hellospider
```

## 相关细节

### 字符编码转换

Go 语言默认字符编码为 ```UTF-8```，在解析网页时为了避免乱码需要判断该网页的字符编码并进行转换。目前发现包含字符编码的内容如下：

* HTTP 响应头中的 ```Content-Type``` ，如：```text/html; charset=gbk```
* HTML 头部的 meta 标签，如：```<meta charset='gbk'>```
  , ```<meta http-equiv='Content-Type' content='text/html; charset=gb2312'>```

如果在响应头中能够判断字符编码，则直接将 ```Reader``` 进行转换即可。否则需要先解析 HTML，判断后再进行转换。后者需要储存响应的数据再次构造 ```Reader``` 以提供重复读取。

### 优先级

每抓取一个网页，往往会产生几十个新的 URL，如果不为 URL 设定优先级策略，或许很难达到预期的结果。

#### a. 无优先级

URL 无优先级时，根据队列先进先出的特性，搜索将会是广度优先。 此时的广度优先并非指的是网站广度，而是相当于树的层次遍历。 搜寻结果根据网页的链接相关性有不同结果。

如果网页内的大部分链接都指向某一个网站，假设是它本站，那么抓取的网页大部分将是此网站的内容或者是与其密切相关的其他站点的内容。

*此策略适合挖掘某个站点*

#### b. 根据 URL 长度计算优先级

一个 URL 越长，通常可能表示它处在站点越深的地方。在某些搜索场景下，它可能并不那么重要。反过来，一个 URL 越短那么它就处于站点越浅的地方，对于想要发现更多站点的搜索场景来说更好。

简单通过一万条 URL 的采样发现 URL 的平均全长在 70 左右，仅 URI 路径的平均长度在 35 左右。

简单设计一个优先级函数：f(x)=(e^((-(x-340))/50))/100

![fx](docs/img/fx.svg?raw=true)

函数图像：

![函数图像](docs/img/fx.jpg?raw=true)

经过实验发现，使用此优先级函数进行的任务能够抓取更多的网站（域名）。

*此策略适合快速发现更多站点*

## 优化计划

* 入队前筛选有效 URL ，避免消息队列臃肿。 √
* 使用协程并发抓取。 √
* 消息队列 Qos。 √
* HTTP 连接复用。 ×
* 维护本地布隆过滤器，避免频繁查询 RedisBloom。×
* HTTP 重定向时，将重定向过程中的 URL 也加入过滤器。 ×

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

### 可视化分析

使用 ```Kibana``` 进行可视化分析。

[Kibana](https://github.com/elastic/kibana)

