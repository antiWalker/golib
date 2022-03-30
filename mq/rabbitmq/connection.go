package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/antiWalker/golib/common"
	"time"
)

type Connection struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	addr          string
	exchange      string
	exchangeType  string
	durable       bool // 消息 是否持久化，rabbit重启之后还会存在，放到磁盘中
	connected     chan interface{}
	quit          chan struct{}
}

func NewConnection(addr string, exchange string, exchangeType string, durable bool) *Connection {
	c := &Connection{
		addr:         addr,
		exchange:     exchange,
		exchangeType: exchangeType,
		durable:      durable,
		connected:    make(chan interface{}, 10),
		quit:         make(chan struct{}),
	}
	return c
}

func (c *Connection) Connect() error {
	var err error
	if c.conn, err = amqp.Dial(c.addr); err != nil {
		common.ErrorLogger.Fatal("rabbitMQ dial to: ", c.addr, "failed:", err)
		return err
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		common.ErrorLogger.Fatal(c.addr, "create channel failed: ", err)
		_ = c.conn.Close()
		return err
	}

	if err = c.channel.ExchangeDeclare(c.exchange, c.exchangeType, c.durable, false, false, false, nil); err != nil {
		common.ErrorLogger.Fatal(c.addr, "declare exchange failed: ", err)
		_ = c.conn.Close()
		return err
	}

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error)) // channel.NotifyClose和connection.NotifyClose可以接收到错误消息，那就以此为重连的触发器
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))
	common.InfoLogger.Info("rabbitmq connect success")
	c.connected <- true
	go c.ReConnect()
	return err
}

func (c *Connection) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				common.ErrorLogger.Info("rabbitmq consumer - connection NotifyClose: ", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				common.ErrorLogger.Info("rabbitmq consumer - channel NotifyClose: ", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel("", true); err != nil {
				common.ErrorLogger.Info("rabbitmq consumer - channel cancel failed: ", err)
			}

			if err := c.conn.Close(); err != nil {
				common.ErrorLogger.Info("rabbitmq consumer - channel cancel failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}

		for {
			select {
			case <-c.quit:
				return
			default:
				time.Sleep(time.Millisecond*500)
				common.InfoLogger.Info("rabbitmq consumer - reconnect")
				if err := c.Connect(); err != nil {
					common.ErrorLogger.Info("rabbitmq consumer - failCheck: ", err)
					continue
				}
				return
			}
		}
	}
}
