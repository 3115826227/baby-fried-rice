package handle

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/module/file/config"
	"baby-fried-rice/internal/pkg/module/file/db"
	"baby-fried-rice/internal/pkg/module/file/log"
	"encoding/json"
)

func InitBackend() {
	conf := config.GetConfig()
	if err := NewConsume(conf.MessageQueue.ConsumeTopics.DeleteFile, runDeleteFileConsume); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

func NewConsume(consume models.TopicConsume, handle func(mq interfaces.MQ)) (err error) {
	conf := config.GetConfig()
	consumeMQ := nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	if err = consumeMQ.NewConsumer(consume.Topic, consume.Channel); err != nil {
		return
	}
	go handle(consumeMQ)
	return
}

func runDeleteFileConsume(consumeMQ interfaces.MQ) {
	for {
		value, err := consumeMQ.Consume()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var info models.DeleteFileMessageQueueInfo
		if err = json.Unmarshal([]byte(value), &info); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		switch info.FileValueType {
		case models.FileId:
			if err = db.GetDB().GetDB().Where("id = ?", info.FileValue).Delete(&tables.File{}).Error; err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		case models.FileDownUrl:
			if err = db.GetDB().GetDB().Where("down_url = ?", info.FileValue).Delete(&tables.File{}).Error; err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		}
	}
}
