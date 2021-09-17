package handle

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/module/manage/config"
	"baby-fried-rice/internal/pkg/module/manage/log"
)

var (
	mq interfaces.MQ
)

func InitBackend() {
	if err := NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

func NewProducer() error {
	conf := config.GetConfig()
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	return mq.NewProducer()
}
