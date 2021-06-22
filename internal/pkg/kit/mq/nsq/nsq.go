package nsq

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"errors"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
	Log "log"
)

var (
	consumerChan = make(chan string, 2000)
)

type nsqMQ struct {
	addr     string
	config   *nsq.Config
	producer *nsq.Producer
	consumer *nsq.Consumer
}

func InitNSQMQ(addr string) interfaces.MQ {
	return &nsqMQ{
		addr:   addr,
		config: nsq.NewConfig(),
	}
}

func (mq *nsqMQ) NewProducer() (err error) {
	if mq.producer != nil {
		return errors.New("producer is exist")
	}
	mq.producer, err = nsq.NewProducer(mq.addr, mq.config)
	if err != nil {
		return
	}
	mq.producer.SetLogger(Log.New(ioutil.Discard, "", Log.LstdFlags), nsq.LogLevelInfo)
	err = mq.producer.Ping()
	if err != nil {
		return
	}
	return
}

func (mq *nsqMQ) Send(topic, value string) error {
	return mq.producer.Publish(topic, []byte(value))
}

func (mq *nsqMQ) NewConsumer(topic, channel string) (err error) {
	if mq.consumer != nil {
		return errors.New("consumer is exist")
	}
	consumerChan = make(chan string, 2000)
	mq.config.MaxInFlight = 10
	mq.consumer, err = nsq.NewConsumer(topic, channel, mq.config)
	if err != nil {
		return
	}
	mq.consumer.SetLogger(Log.New(ioutil.Discard, "", Log.LstdFlags), nsq.LogLevelInfo)
	mq.consumer.AddHandler(&nsqMQ{})
	if err = mq.consumer.ConnectToNSQD(mq.addr); err != nil {
		return
	}
	return
}

func (mq *nsqMQ) Consume() (value string, err error) {
	var ok bool
	value, ok = <-consumerChan
	if !ok {
		err = errors.New("failed to consume")
	}
	return
}

func (mq *nsqMQ) HandleMessage(msg *nsq.Message) (err error) {
	consumerChan <- string(msg.Body)
	return
}
