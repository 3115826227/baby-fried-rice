package handle

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/module/im/config"
	"baby-fried-rice/internal/pkg/module/im/log"
)

var (
	mq    interfaces.MQ
	topic string
)

func Init() {
	conf := config.GetConfig()
	topic = conf.NSQ.Topic
	mq = nsq.InitNSQMQ(conf.NSQ.Addr)
	if err := mq.NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
