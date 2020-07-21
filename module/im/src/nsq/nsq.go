package nsq

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
	"io/ioutil"
	Log "log"
)

var producer *nsq.Producer
var consumer *nsq.Consumer

func init() {
	cfg := nsq.NewConfig()

	var err error
	producer, err = nsq.NewProducer(config.Config.Nsq, cfg)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	producer.SetLogger(Log.New(ioutil.Discard, "", Log.LstdFlags), nsq.LogLevelInfo)
	err = producer.Ping()
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	cfg.MaxInFlight = 10
	consumer, err = nsq.NewConsumer(config.ChatTopic, config.ConsumerChatChannel, cfg)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	consumer.SetLogger(Log.New(ioutil.Discard, "", Log.LstdFlags), nsq.LogLevelInfo)
	consumer.AddHandler(&HandlerNSQ{})
	if err := consumer.ConnectToNSQD(config.Config.Nsq); err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}

type HandlerNSQ struct {
}

var ChatCh = make(chan model.FriendChatMessageReq, 5000)

func (handle *HandlerNSQ) HandleMessage(msg *nsq.Message) (err error) {
	var message model.FriendChatMessageReq
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	ChatCh <- message
	log.Logger.Info("receive", zap.String("addr", msg.NSQDAddress), zap.String("msg", string(msg.Body)))
	return
}

func Send(message, key, topic string) (err error) {
	err = producer.Publish(topic, []byte(message))
	if err != nil {
		log.Logger.Warn("send failed: " + err.Error())
		return
	}
	log.Logger.Info("send successful: " + message)
	return
}
