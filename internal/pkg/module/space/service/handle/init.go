package handle

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/module/space/config"
	"baby-fried-rice/internal/pkg/module/space/log"
)

var (
	mq    interfaces.MQ
	topic string
)

func Init() {
	conf := config.GetConfig()
	topic = conf.MessageQueue.PublishTopics.WebsocketNotify
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	if err := mq.NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
