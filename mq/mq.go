package mq

import (
	"github.com/NingziSlay/pkg/log"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"sync"
)

type ExchangeKind string

const (
	ExchangeTopic    ExchangeKind = "topic"
	ExchangeTFanout  ExchangeKind = "fanout"
	ExchangeTDirect  ExchangeKind = "direct"
	ExchangeTHeaders ExchangeKind = "headers"
)

type Config struct {
	Addr          string
	Exchange      string
	ExchangeType  ExchangeKind // topic, direct, etc
	Queue         string
	RoutingKey    string
	ConsumerTag   string
	PrefetchCount int
	PrefetchSize  int
	ExchangeArgs  amqp.Table
	QueueArgs     amqp.Table
	QueueBindArgs amqp.Table
}

type mq struct {
	once sync.Once

	conn    *amqp.Connection
	channel *amqp.Channel
	config  *Config
	log     zerolog.Logger
}

func newMq(config *Config) *mq {
	return &mq{
		once:   sync.Once{},
		config: config,
		log:    log.GetLogger(),
	}
}

// stop 关闭 consumer
func (q *mq) stop() {
	if !q.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := q.channel.Cancel(q.config.ConsumerTag, true); err != nil {
			q.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
		}

		if err := q.conn.Close(); err != nil {
			q.log.Warn().Err(err).Msg("rabbitmq consumer - connection close failed")
		}
	}
}

// init exchange、queue、queue bind 都做了冗余的声明操作，为了防止发送的消息
// 在 mq server 里匹配不到对应的 queue
func (q *mq) init() (err error) {
	q.once.Do(func() {
		if q.conn, err = amqp.Dial(q.config.Addr); err != nil {
			return
		}

		if q.channel, err = q.conn.Channel(); err != nil {
			q.stop()
			return
		}

		if err = q.channel.ExchangeDeclare(q.config.Exchange, string(q.config.ExchangeType), true, false, false, false, q.config.ExchangeArgs); err != nil {
			q.stop()
			return
		}

		if _, err = q.channel.QueueDeclare(q.config.Queue, true, false, false, false, q.config.QueueArgs); err != nil {
			q.stop()
			return
		}

		_ = q.channel.Qos(q.config.PrefetchCount, q.config.PrefetchSize, true)

		if err = q.channel.QueueBind(q.config.Queue, q.config.RoutingKey, q.config.Exchange, false, q.config.QueueBindArgs); err != nil {
			q.stop()
			return
		}
		return
	})
	return
}
