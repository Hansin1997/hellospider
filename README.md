# Hello Spider ð·

åºäº ```Go``` è¯­è¨çåå¸å¼ç½é¡µç¬è«ã

*ç¬åæ°æ®åå¯ä½¿ç¨ 
[Hello Search](https://github.com/Hansin1997/hellosearch)
è¿è¡æç´¢ã*

## ç®ä»

> æ­¤é¡¹ç®åè¡·å¨äºå­¦ä¹  ```Go``` è¯­è¨ä»¥å ```Elasticsearch``` ã

### åºæ¬åç

1. å°ç§å­ URL å¥éã
2. ä»éåè·åå¾æå URLã
3. å¤æ­æ­¤ URL æ¯å¦ææãï¼åæ¬æ£æ¥ URL æ¯å¦å·²ç»è¢«æåè¿ï¼
4. æåç½é¡µä¿¡æ¯ï¼å¹¶è®°å½æ­¤ URL å·²ç»è¢«æåã
5. å¨å­ç½é¡µæè¦ï¼å¹¶å°æ­¤ç½é¡µçè¶é¾æ¥ï¼URLï¼è¿è¡å¥éã
6. æ§è¡æ­¥éª¤ ``` 2```ã

### å®ç°æ¦è¦

* ä½¿ç¨ ```RabbitMQ``` æ¶æ¯éåå®ç°åå¸å¼ä¸ç URL éåçæä¹åãä¼åçº§ç­ã
* ä½¿ç¨ ```RedisBloom``` çå¸éè¿æ»¤å¨å®ç°åå¸å¼ä¸ç URL å»éã
* ä½¿ç¨ ```goquery``` è§£æ HTMLã
* ä½¿ç¨ ```Elasticsearch``` å¨å­æ°æ®ã

å¦å¤å¯ä»¥ä½¿ç¨ ```Kibana``` å¯¹æåç¶åµè¿è¡å®æ¶çå¯è§ååæã

### Kibana å¯è§å

![Kibanaå¯è§å](docs/img/kibana.png?raw=true)

## ç¸å³ç¯å¢

### æå»ºç¯å¢

* ```go version go1.16.3```

### è¿è¡ç¯å¢

* ```RabbitMQ 3.8.12```
* ```RedisBloom latest```

## ä½¿ç¨æ¹æ³

### å½ä»¤è¡åæ°

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

* ```-config``` éç½®æä»¶è·¯å¾
* ```-namespace``` å½åç©ºé´ï¼åºåä¸åä»»å¡ï¼
* ```-priority``` ä¼åçº§ç­ç¥ï¼0-9 è¡¨ç¤ºä¼åçº§ä¸ºå¸¸æ°ï¼url-len è¡¨ç¤ºæ ¹æ® URL é¿åº¦è®¡ç®ä¼åçº§ï¼è¶ç­è¶ä¼åï¼ï¼path-len è¡¨ç¤ºæ ¹æ® URL è·¯å¾é¿åº¦è®¡ç®ä¼åçº§ï¼è¶ç­è¶ä¼åï¼ã
* ```-reset``` å¼å§åæ¯å¦éç½®å½åç©ºé´
* ```-seed``` æ¿æ¢éç½®æä»¶ä¸­ç URL ç§å­ï¼è±æéå·åéï¼

### éç½®æä»¶

ä¿®æ¹éç½®æä»¶ ```config.json``` ä¸­åæå¡çå°åãç«¯å£ä»¥åç¨æ·ååå¯ç ç­ã

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
  "rules": {
    "allows": [
      ".*"
    ],
    "forbid": null
  },
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

### è¿è¡

#### æºä»£ç 

```bash
go run .
```

### äºè¿å¶å¯æ§è¡æä»¶

```bash
./hellospider
```

## ç¸å³ç»è

### å­ç¬¦ç¼ç è½¬æ¢

Go è¯­è¨é»è®¤å­ç¬¦ç¼ç ä¸º ```UTF-8```ï¼å¨è§£æç½é¡µæ¶ä¸ºäºé¿åä¹±ç éè¦å¤æ­è¯¥ç½é¡µçå­ç¬¦ç¼ç å¹¶è¿è¡è½¬æ¢ãç®ååç°åå«å­ç¬¦ç¼ç çåå®¹å¦ä¸ï¼

* HTTP ååºå¤´ä¸­ç ```Content-Type``` ï¼å¦ï¼```text/html; charset=gbk```
* HTML å¤´é¨ç meta æ ç­¾ï¼å¦ï¼```<meta charset='gbk'>```
  , ```<meta http-equiv='Content-Type' content='text/html; charset=gb2312'>```

å¦æå¨ååºå¤´ä¸­è½å¤å¤æ­å­ç¬¦ç¼ç ï¼åç´æ¥å° ```Reader``` è¿è¡è½¬æ¢å³å¯ãå¦åéè¦åè§£æ HTMLï¼å¤æ­ååè¿è¡è½¬æ¢ãåèéè¦å¨å­ååºçæ°æ®åæ¬¡æé  ```Reader``` ä»¥æä¾éå¤è¯»åã

### ä¼åçº§

æ¯æåä¸ä¸ªç½é¡µï¼å¾å¾ä¼äº§çå åä¸ªæ°ç URLï¼å¦æä¸ä¸º URL è®¾å®ä¼åçº§ç­ç¥ï¼æè®¸å¾é¾è¾¾å°é¢æçç»æã

#### a. æ ä¼åçº§

URL æ ä¼åçº§æ¶ï¼æ ¹æ®éååè¿ååºçç¹æ§ï¼æç´¢å°ä¼æ¯å¹¿åº¦ä¼åã æ­¤æ¶çå¹¿åº¦ä¼åå¹¶éæçæ¯ç½ç«å¹¿åº¦ï¼èæ¯ç¸å½äºæ çå±æ¬¡éåã æå¯»ç»ææ ¹æ®ç½é¡µçé¾æ¥ç¸å³æ§æä¸åç»æã

å¦æç½é¡µåçå¤§é¨åé¾æ¥é½æåæä¸ä¸ªç½ç«ï¼åè®¾æ¯å®æ¬ç«ï¼é£ä¹æåçç½é¡µå¤§é¨åå°æ¯æ­¤ç½ç«çåå®¹æèæ¯ä¸å¶å¯åç¸å³çå¶ä»ç«ç¹çåå®¹ã

*æ­¤ç­ç¥éåæææä¸ªç«ç¹*

#### b. æ ¹æ® URL é¿åº¦è®¡ç®ä¼åçº§

ä¸ä¸ª URL è¶é¿ï¼éå¸¸å¯è½è¡¨ç¤ºå®å¤å¨ç«ç¹è¶æ·±çå°æ¹ãå¨æäºæç´¢åºæ¯ä¸ï¼å®å¯è½å¹¶ä¸é£ä¹éè¦ãåè¿æ¥ï¼ä¸ä¸ª URL è¶ç­é£ä¹å®å°±å¤äºç«ç¹è¶æµçå°æ¹ï¼å¯¹äºæ³è¦åç°æ´å¤ç«ç¹çæç´¢åºæ¯æ¥è¯´æ´å¥½ã

ç®åéè¿ä¸ä¸æ¡ URL çéæ ·åç° URL çå¹³åå¨é¿å¨ 70 å·¦å³ï¼ä» URI è·¯å¾çå¹³åé¿åº¦å¨ 35 å·¦å³ã

ç®åè®¾è®¡ä¸ä¸ªä¼åçº§å½æ°ï¼f(x)=(e^((-(x-340))/50))/100

![fx](docs/img/fx.svg?raw=true)

å½æ°å¾åï¼

![å½æ°å¾å](docs/img/fx.jpg?raw=true)

ç»è¿å®éªåç°ï¼ä½¿ç¨æ­¤ä¼åçº§å½æ°è¿è¡çä»»å¡è½å¤æåæ´å¤çç½ç«ï¼ååï¼ã

*æ­¤ç­ç¥éåå¿«éåç°æ´å¤ç«ç¹*

## ä¼åè®¡å

* å¥éåç­éææ URL ï¼é¿åæ¶æ¯éåèè¿ã â
* ä½¿ç¨åç¨å¹¶åæåã â
* æ¶æ¯éå Qosã â
* HTTP è¿æ¥å¤ç¨ã Ã
* ç»´æ¤æ¬å°å¸éè¿æ»¤å¨ï¼é¿åé¢ç¹æ¥è¯¢ RedisBloomãÃ
* HTTP éå®åæ¶ï¼å°éå®åè¿ç¨ä¸­ç URL ä¹å å¥è¿æ»¤å¨ã Ã

## ç¸å³ææ¯

### å¸éè¿æ»¤å¨

```RedisBloom``` ç¨äºåå¸å¼ä¸ç ```URL``` å»éã

[RedisBloom](https://github.com/RedisBloom/RedisBloom)

[redisbloom-go](https://github.com/RedisBloom/redisbloom-go)

### æ¶æ¯éå

ä½¿ç¨ ```RabbitMQ``` å­æ¾å¾æåç ```URL```ï¼å¹¶ä»éåè·å ```URL``` è¿è¡æ¶è´¹ã

[RabbitMQ Server](https://github.com/rabbitmq/rabbitmq-server)

[Go RabbitMQ Client Library](https://github.com/streadway/amqp)

### æ°æ®å­å¨

ä½¿ç¨ ```Elasticsearch``` å­å¨ç½é¡µæ°æ®ã

[Elasticsearch](https://github.com/elastic/elasticsearch)

[go-elasticsearch](https://github.com/elastic/go-elasticsearch)

### å¯è§ååæ

ä½¿ç¨ ```Kibana``` è¿è¡å¯è§ååæã

[Kibana](https://github.com/elastic/kibana)

