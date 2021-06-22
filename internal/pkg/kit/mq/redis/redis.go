package redis

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

type redisMQ struct {
	addr         string
	password     string
	db           int
	producer     *redis.Client
	consumer     *redis.Client
	consumerChan chan string
}

func InitRedisMQ(addr, password string, db int) interfaces.MQ {
	return &redisMQ{
		addr:         addr,
		password:     password,
		db:           db,
		consumerChan: make(chan string, 2000),
	}
}

func (mq *redisMQ) NewProducer() (err error) {
	if mq.producer != nil {
		return errors.New("producer is exist")
	}
	mq.producer = redis.NewClient(&redis.Options{
		Addr:     mq.addr,
		Password: mq.password,
		PoolSize: 20,
		DB:       mq.db,
	})
	if err = mq.producer.Ping().Err(); err != nil {
		return
	}
	return
}

func (mq *redisMQ) Send(topic, value string) error {
	return mq.producer.LPush(topic, value).Err()
}

func (mq *redisMQ) NewConsumer(topic, channel string) (err error) {
	if mq.consumer != nil {
		return errors.New("consumer is exist")
	}
	mq.consumer = redis.NewClient(&redis.Options{
		Addr:     mq.addr,
		Password: mq.password,
		PoolSize: 20,
		DB:       mq.db,
	})
	if err = mq.consumer.Ping().Err(); err != nil {
		return
	}
	for {
		var results []string
		results, err = mq.consumer.BLPop(1*time.Second, topic).Result()
		if err != nil {
			return
		}
		for _, res := range results {
			mq.consumerChan <- res
		}
	}
}

func (mq *redisMQ) Consume() (value string, err error) {
	var ok bool
	value, ok = <-mq.consumerChan
	if !ok {
		err = errors.New("failed to consume")
	}
	return
}
