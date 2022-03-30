package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/nacos"
	"time"
)

type Publisher struct {
	conn     *Connection
	exchange string
	routingKey string
	exchangeType string
}

func GetConfig() *Publisher {
	defer func() {
		if err:= recover();err != nil{
			common.ErrorLogger.Info("rabbitmq NewPublisher Error:", err)
		}
	}()
	config := nacos.GetConfigMap()
	RabbitMqUserName := config["rabbitmq.username"]
	RabbitMqPassWord := config["rabbitmq.password"]
	RabbitMqHost := config["rabbitmq.host"]
	RabbitMqExchange := config["rabbitmq.exchange"]
	RabbitMqExchangeType := config["rabbitmq.exchangeType"]
	RabbitMqRoutingKey := config["rabbitmq.routingKey"]
	// 构造链接地址
	addr := "amqp://" + RabbitMqUserName + ":" + RabbitMqPassWord + "@" + RabbitMqHost
	return NewPublisher(addr, RabbitMqExchange, RabbitMqExchangeType, RabbitMqRoutingKey)
}

func NewPublisher(addr, RabbitMqExchange, RabbitMqExchangeType,RabbitMqRoutingKey string) *Publisher {

	c := &Publisher{
		conn: NewConnection(addr, RabbitMqExchange, RabbitMqExchangeType, true),
		exchange: RabbitMqExchange,
		routingKey: RabbitMqRoutingKey,
	}

	_ = c.conn.Connect()
	return c
}

func (p *Publisher) Send(msg string) {
	for {
	clear:
		for {
			select {
			case <-p.conn.connected:
			default:
				break clear
			}
		}
		if err := p.conn.channel.Publish(
			p.conn.exchange, // exchange
			p.routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				DeliveryMode: amqp.Transient,
				Body:         []byte(msg),
			}); err != nil {
			common.InfoLogger.Info("rabbitmq publish - failCheck: ", err)
			time.Sleep(time.Millisecond*100)
			continue
		}
		break
	}
}
