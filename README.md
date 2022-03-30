# 基础包

所有的基础公共包都放在这个项目下，统一管理。

### cache 缓存相关的公共包

* redis/redis 包，配置信息从nacos中获取，格式如下：
```
redis.Host = 127.0.0.1:6379
redis.Password = 123456
redis.DB = 0
redis.DialTimeout = 10
redis.ReadTimeout = 30
redis.WriteTimeout = 30
redis.PoolSize = 100
redis.PoolTimeout = 30
```

### common 公共方法提取的公共包

* common/ip_utils.go   : ip地址公共方法

* common/log.go   :  集成uber zap 日志框架，四种文件记录级别

```
var InfoLogger *zap.SugaredLogger
var DebugLogger *zap.SugaredLogger
var ErrorLogger *zap.SugaredLogger
var WarnLogger *zap.SugaredLogger
```

* common/reflect_util.go   : 反射相关的公共方法

### db 数据源注册的公共包

* orm/orm.go : orm支持的数据库类型

```
mysql | sqlite | oracle | pgsql | TiDB
```

**支持多数据源的配置文件从nacos中获取，格式为：**

```
mysqlkeymultiple=default::QRQM,slave::WEBAPP
mysqlkey = QRQM
QRQM.Host = 127.0.0.1
QRQM.Port = 3306
QRQM.DbName = nicetuan_recommend_qrqm
QRQM.UserName = root
QRQM.Password = 123456
QRQM.DriverName = mysql
WORD.SetMaxIdleConns = 50
WORD.SetMaxOpenConns = 50

slave = WEBAPP
WEBAPP.Host = 127.0.0.1
WEBAPP.Port = 3306
WEBAPP.DbName = nicetuan_weapp_db_new
WEBAPP.UserName = root
WEBAPP.Password = 123456
WEBAPP.DriverName = mysql
```

### monitor 监控报警相关的公共包

**dingding/dingding.go : 钉钉报警，配置信息从nacos中获取，格式如下**
```
dingding.clientId=recommend
dingding.keyWords=推荐预警
dingding.security=123
dingding.webHook=https://oapi.dingtalk.com/robot/send?access_token=
```

### mq 消息队列相关的公共包

* kafka/kafka.go

**生产者和消费者都存放在此处，配置文件从nacos获取，格式如下：**
```
kafka.producer.bootstrapServers = 127.0.0.1:9092
kafka.producer.topic = test_rec_bigdata_recommend_result

kafka.consumer.bootstrapServers = 127.0.0.1:9092
kafka.consumer.topics = rec_bigdata_user_action_accum
kafka.consumer.groupId = recommend_realtime_group_golang_local
kafka.consumer.count = 10
```
* rabbitmq/publisher.go
**配置文件从nacos获取，格式如下：**
```
rabbitmq.username = admin
rabbitmq.password = admin
rabbitmq.host = 127.0.0.1:5672
rabbitmq.exchange = sht.sco.alarmMessage
rabbitmq.exchangeType = topic
rabbitmq.routingKey = scoAlarmMessage

```
### nacos 配置中心相关的公共包

**配置中心使用nacos，所有的配置信息都会存放在nacos上。nacos配置文件存在程序的当前目录下： ./conf/nacos.properties**

```
ip = 127.0.0.1
port = 8847
namespaceId = test
dataId = slp_recommend
group = recommend
logDir = ./logs/nacos/log
cacheDir = ./conf
rotateTime = 1h
logLevel = debug
timeoutMs = 5000
maxAge = 3
```

