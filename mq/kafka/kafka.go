package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/nacos"
	"strconv"
	"strings"
	"time"
)

// SendMessageWithTopicName 发送消息，指定topicName
func SendMessageWithTopicName(topicName, key, value string) (partition int32, offset int64, err error) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topicName
	msg.Value = sarama.StringEncoder(value)
	msg.Key = sarama.StringEncoder(key)
	return SendProducerMessage(msg)
}

func SendMessage(key string, value string) (partition int32, offset int64, err error) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = GetProducerTopic()
	msg.Value = sarama.StringEncoder(value)
	msg.Key = sarama.StringEncoder(key)
	return SendProducerMessage(msg)
}

func SendProducerMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	// 发送消息
	client, err := GetKafkaClient()
	if err != nil {
		common.ErrorLogger.Infof("Error get kafka client : %v \n", err)
	}
	defer client.Close()

	return client.SendMessage(msg)
}

func makeProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	return config
}

func GetKafkaClient() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(GetProducerBootstrapServers(), makeProducerConfig())
}

func GetProducerBootstrapServers() []string {
	servers := strings.Split(nacos.GetConfigMap()["kafka.producer.bootstrapServers"], ",")
	for _, server := range servers {
		bootstrapServers := strings.Split(server, ",")
		return bootstrapServers
	}
	return nil
}

// GetProducerTopicName 获取生产者的topicName
func GetProducerTopicName(nacosKey string) string {
	return nacos.GetConfigMap()[nacosKey]
}

// GetProducerTopic 获取生产者的topic
func GetProducerTopic() string {
	return nacos.GetConfigMap()["kafka.producer.topic"]
}

func makeConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Minute
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	return config
}

func GetConsumerGroup() (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(GetConsumerBootstrapServers(), GetConsumerGroupId(), makeConsumerConfig())
}

func GetConsumerBootstrapServers() []string {
	servers := strings.Split(nacos.GetConfigMap()["kafka.consumer.bootstrapServers"], ",")
	for _, server := range servers {
		bootstrapServers := strings.Split(server, ",")
		return bootstrapServers
	}
	return nil
}

// GetConsumerGroupId 获取消费者组
func GetConsumerGroupId() string {
	return nacos.GetConfigMap()["kafka.consumer.groupId"]
}

// GetConsumerTopics 获取消费者的top集合
func GetConsumerTopics() []string {
	return strings.Split(nacos.GetConfigMap()["kafka.consumer.topics"], ",")
}

// GetConsumerCount 获取消费者的数量
func GetConsumerCount() (int, error) {
	return strconv.Atoi(nacos.GetConfigMap()["kafka.consumer.count"])
}
