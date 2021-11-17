package kafka

import (
	"github.com/Shopify/sarama"
	"time"
)

type kafkaMQ struct {
	addr     string
	producer sarama.AsyncProducer

	newConsumerNotify chan string
	consumer          sarama.Consumer
	consumerChan      chan string
}

func (mq *kafkaMQ) run() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	consumer, err := sarama.NewConsumer([]string{mq.addr}, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	for {
		select {
		case <-mq.newConsumerNotify:
			// 新增了消费者
		}
	}
}

func (mq *kafkaMQ) NewProducer() error {
	if mq.producer != nil {
		return nil
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2
	producer, err := sarama.NewAsyncProducer([]string{mq.addr}, config)
	if err != nil {
		return err
	}
	mq.producer = producer
	return nil
}

func (mq *kafkaMQ) Send(topic, value string) error {
	mq.producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(value),
		Timestamp: time.Now(),
	}
	return nil
}

func (mq *kafkaMQ) NewConsumer(topic, partition string) (err error) {
	return
}

func (mq *kafkaMQ) Consume() (value string, err error) {
	return
}
