package handle

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/module/shop/config"
	"baby-fried-rice/internal/pkg/module/shop/log"
)

var (
	mq interfaces.MQ
)

func Init() {
	conf := config.GetConfig()
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	if err := mq.NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
