package redis

import (
	"github.com/go-redis/redis"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/nacos"
	"os"
	"strconv"
	"time"
)

var (
	RedisCli *redis.Client
)

func getRedisInstance() *redis.Client {
	defer func() {
		e := recover()
		if e != nil {
			common.ErrorLogger.Infof("redis recover: %s", e)
		}
	}()
	if RedisCli != nil {
		_, err := RedisCli.Ping().Result()
		if err == nil {
			return RedisCli
		}
	}

	config := nacos.GetConfigMap()
	redisHost := config["redis.Host"]
	redisPassword := config["redis.Password"]

	envk8s := os.Getenv("RUNTIME_TYPE")
	var RedisAddr string
	if envk8s == "product" {
		RedisAddr = redisHost
	} else {
		RedisAddr = redisHost
	}
	redisDB, _ := strconv.Atoi(config["redis.DB"])
	redisDialTimeout, _ := strconv.Atoi(config["redis.DialTimeout"])
	redisReadTimeout, _ := strconv.Atoi(config["redis.ReadTimeout"])
	redisWriteTimeout, _ := strconv.Atoi(config["redis.WriteTimeout"])
	redisPoolSize, _ := strconv.Atoi(config["redis.PoolSize"])
	redisPoolTimeout, _ := strconv.Atoi(config["redis.PoolTimeout"])

	RedisCli = redis.NewClient(&redis.Options{
		Addr:         RedisAddr,
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  time.Duration(redisDialTimeout) * time.Second,
		ReadTimeout:  time.Duration(redisReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisWriteTimeout) * time.Second,
		PoolSize:     redisPoolSize,
		PoolTimeout:  time.Duration(redisPoolTimeout) * time.Second,
	})
	if _, err := RedisCli.Ping().Result(); err != nil {
		common.ErrorLogger.Infof("write redis error: %v", err)
	}
	_, err := RedisCli.Ping().Result()
	if err != nil {
		return nil
	}
	return RedisCli
}

func RedisGet(key string) string {
	client := getRedisInstance()

	data, err := client.Get(key).Result()

	if err != nil {
		data = ""
	}

	return data
}

func RedisMGet(key []string) (data []interface{}, err error) {
	client := getRedisInstance()

	data, err = client.MGet(key...).Result()
	if err != nil {
		common.ErrorLogger.Infof("MGet redis error: %s", err)
	}

	return data, err
}

func RedisSet(key string, value interface{}) {

	RedisSetEx(key, value, 0)
}
func RedisSetEx(key string, value interface{}, expiration time.Duration) {
	client := getRedisInstance()
	client.Set(key, value, expiration)
}

func RedisIncr(key string) {
	RedisIncrEx(key, 0)
}

func RedisIncrEx(key string, expiration time.Duration) int64 {
	client := getRedisInstance()
	incr := client.Incr(key)
	//第一次设置key，需要添加过期时间
	if expiration > 0 && incr.Val() == 1 {
		//设置过期时间
		client.Expire(key, expiration)
	}
	return incr.Val()
}

/**
 * @Description:	将哈希表 hash 中域 field 的值设置为 value 。
 * @param key
 * @param field
 * @param value
 * @return bool
 */
func RedisHSet(key string, field string, value interface{}) bool {
	result, err := getRedisInstance().HSet(key, field, value).Result()
	if err != nil {
		result = false
	}
	return result
}

/**
 * @Description:当且仅当域 field 尚未存在于哈希表的情况下， 将它的值设置为 value 。如果给定域已经存在于哈希表当中， 那么命令将放弃执行设置操作。
 * @param key
 * @param field
 * @param value
 * @return bool
 */
func RedisHSetNx(key string, field string, value interface{}) bool {
	result, err := getRedisInstance().HSetNX(key, field, value).Result()
	if err != nil {
		result = false
	}
	return result
}

/**
 * @Description:同时将多个 field-value (域-值)对设置到哈希表 key 中。此命令会覆盖哈希表中已存在的域。
 * @param key
 * @param fields
 * @return string
 */
func RedisHMSet(key string, fields map[string]interface{}) string {
	result, err := getRedisInstance().HMSet(key, fields).Result()
	if err != nil {
		return ""
	}
	return result
}

/**
 * @Description:	返回哈希表中给定域的值。
 * @param key
 * @param field
 * @return string
 */
func RedisHGet(key string, field string) string {
	result, err := getRedisInstance().HGet(key, field).Result()
	if err != nil {
		result = ""
	}
	return result
}

/**
 * @Description:	返回哈希表 key 中，一个或多个给定域的值。
 * @param key
 * @param fields
 * @return data		一个包含多个给定域的关联值的表，表值的排列顺序和给定域参数的请求顺序一样。
 * @return err
 */
func RedisHMGet(key string, fields ...string) (data []interface{}, err error) {
	result, err := getRedisInstance().HMGet(key, fields...).Result()
	if err != nil {
		common.ErrorLogger.Infof("HMGet redis error: %s", err)
	}
	return result, err
}

/**
 * @Description:	检查给定域 field 是否存在于哈希表 hash 当中。
 * @param key
 * @param field
 * @return bool
 */
func RedisHExists(key string, field string) bool {
	result, err := getRedisInstance().HExists(key, field).Result()
	if err != nil {
		result = false
	}
	return result
}

/**
 * @Description:	删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
 * @param key
 * @param field
 * @return int64
 */
func RedisHDel(key string, field ...string) int64 {
	result, err := getRedisInstance().HDel(key, field...).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	哈希表中域的数量。
 * @param key
 * @return int64
 */
func RedisHLen(key string) int64 {
	result, err := getRedisInstance().HLen(key).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	返回哈希表 key 中的所有域。
 * @param key
 * @return []string
 */
func RedisHKeys(key string) []string {
	result, err := getRedisInstance().HKeys(key).Result()
	if err != nil {
		result = []string{}
	}
	return result
}

/**
 * @Description: 	返回哈希表 key 中所有域的值。
 * @param key
 * @return []string
 */
func RedisHVals(key string) []string {
	result, err := getRedisInstance().HVals(key).Result()
	if err != nil {
		result = []string{}
	}
	return result
}

/**
 * @Description: 	返回哈希表 key 中所有域的值。
 * @param key
 * @return []string
 */
func RedisHGetAll(key string) map[string]string {
	result, err := getRedisInstance().HGetAll(key).Result()
	if err != nil {
		result = make(map[string]string)
	}
	return result
}

/**
 * @Description: 	将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
 * @param key		key
 * @param member	n 个 member
 * @return int64	被添加到集合中的新元素的数量，不包括被忽略的元素。
 */
func RedisSADD(key string, member ...interface{}) int64 {
	result, err := getRedisInstance().SAdd(key, member).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	返回集合 key 中的所有成员。
 * @param key		key
 * @return []string	返回集合 key 中的所有成员。
 */
func RedisSMembers(key string) []string {
	result, err := getRedisInstance().SMembers(key).Result()
	if err != nil {
		result = []string{}
	}
	return result
}

/**
 * @Description:	判断 member 元素是否集合 key 的成员。
 * @param key
 * @param member
 * @return bool		如果 member 元素是集合的成员，返回 true 。 如果 member 元素不是集合的成员，或 key 不存在，返回 false 。
 */
func RedisSIsMember(key string, member interface{}) bool {
	result, err := getRedisInstance().SIsMember(key, member).Result()
	if err != nil {
		result = false
	}
	return result
}

/**
 * @Description:	移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
 * @param key
 * @param member
 * @return int64	被成功移除的元素的数量，不包括被忽略的元素。
 */
func RedisSRem(key string, member ...interface{}) int64 {
	result, err := getRedisInstance().SRem(key, member).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:		将 member 元素从 source 集合移动到 destination 集合。
 * @param sourceKey
 * @param destinationKey
 * @param member
 * @return bool			true:移动成功;false:移动失败
 */
func RedisSMove(sourceKey string, destinationKey string, member ...interface{}) bool {
	result, err := getRedisInstance().SMove(sourceKey, destinationKey, member).Result()
	if err != nil {
		result = false
	}
	return result
}

/**
 * @Description:	返回集合 key 的基数(集合中元素的数量)。
 * @param key
 * @return int64	集合中元素的数量
 */
func RedisSCard(key string) int64 {
	result, err := getRedisInstance().SCard(key).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	求两个set集合中的并集
 * @param key
 * @return []string	并集成员的列表。
 */
func RedisSUnion(key ...string) []string {
	result, err := getRedisInstance().SUnion(key...).Result()
	if err != nil {
		result = []string{}
	}
	return result
}

/**
 * @Description: 	使用zset结构进行添加数据
 * @param key		key
 * @param score		分数
 * @param member	对象
 * @return int64	添加结果
 */
func RedisZADD(key string, score float64, member string) int64 {
	result, err := getRedisInstance().ZAdd(key, redis.Z{
		Score:  score,
		Member: member,
	}).Result()

	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	通过key 和 member 获取分数
 * @param key		key
 * @param member	member
 * @return float64	获member 成员的 score 值
 */
func RedisZCore(key string, member string) float64 {
	result, err := getRedisInstance().ZScore(key, member).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description: 	返回有序集 key 的数量。
 * @param key		key
 * @return int64	数量
 */
func RedisZCard(key string) int64 {
	result, err := getRedisInstance().ZCard(key).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description: 	返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
 * @param key		key
 * @param mix		最小值
 * @param max		最大值
 * @return int64	score 值在 min 和 max 之间的成员的数量。
 */
func RedisZCount(key, mix, max string) int64 {
	result, err := getRedisInstance().ZCount(key, mix, max).Result()
	if err != nil {
		result = 0
	}
	return result
}

/**
 * @Description:	返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从小到大)来排序。
 * @param key		key
 * @param start		开始分数
 * @param stop		结束分数
 * @return []string	指定区间内，带有 score 值(可选)的有序集成员的列表。
 */
func RedisZRange(key string, start int64, stop int64) []string {
	result, err := getRedisInstance().ZRange(key, start, stop).Result()
	if err != nil {
		result = []string{}
	}
	return result
}

/**
 * @Description:	返回有序集 key 中，指定区间内的成员。其中成员的位置按 score 值递增(从大到小)来排序。
 * @param key		key
 * @param start		开始分数
 * @param stop		结束分数
 * @return []string	指定区间内，带有 score 值(可选)的有序集成员的列表。
 */
func RedisZRevRange(key string, start int64, stop int64) []string {
	result, err := getRedisInstance().ZRevRange(key, start, stop).Result()
	if err != nil {
		result = []string{}
	}
	return result
}
