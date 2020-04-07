package kafka

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	producer sarama.SyncProducer
	consumer sarama.Consumer
)

func init() {
	cfg := sarama.NewConfig()

	cfg.Producer.RequiredAcks = sarama.WaitForAll

	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V0_10_0_0

	var err error
	consumer, err = sarama.NewConsumer(strings.Split(config.Config.Kafka, ","), cfg)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	producer, err = sarama.NewSyncProducer(strings.Split(config.Config.Kafka, ","), cfg)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}

func Send(message, key, topic string) (pid int32, err error) {
	var offset int64
	pid, offset, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	log.Logger.Info("kafka send to:",
		zap.Int32("pid", pid), zap.Int64("offset", offset))
	return
}

var ChatCh = make(chan model.FriendChatMessageReq, 5000)

func ReceiveChat(topic string) {
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		for message := range partitionConsumer.Messages() {
			var msg model.FriendChatMessageReq
			err = json.Unmarshal(message.Value, &msg)
			if err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			ChatCh <- msg
		}
	}
	return
}
