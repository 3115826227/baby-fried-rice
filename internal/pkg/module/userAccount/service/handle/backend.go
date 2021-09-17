package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"encoding/json"
)

var (
	mq            interfaces.MQ
	payFailedChan = make(chan models.UserCoinChangeMQMessage, 10)
)

func InitBackend() {
	conf := config.GetConfig()
	if err := NewProducer(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err := NewConsume(conf.MessageQueue.ConsumeTopics.UserCoin, runUserCoinConsume); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

func NewProducer() error {
	conf := config.GetConfig()
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	return mq.NewProducer()
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

// 处理用户积分变动的消息
func runUserCoinConsume(consumeMQ interfaces.MQ) {
	go handlePayFailed()
	for {
		value, err := consumeMQ.Consume()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var msg models.UserCoinChangeMQMessage
		if err = json.Unmarshal([]byte(value), &msg); err != nil {
			log.Logger.Error(err.Error())
			payFailedChan <- msg
			continue
		}
		var userClient user.DaoUserClient
		userClient, err = grpc.GetUserClient()
		if err != nil {
			log.Logger.Error(err.Error())
			payFailedChan <- msg
			continue
		}
		var reqCoin = &user.ReqUserCoinLogAddDao{
			AccountId: msg.AccountId,
			Coin:      msg.Coin,
			CoinType:  msg.CoinType,
		}
		if _, err = userClient.UserCoinLogAddDao(context.Background(), reqCoin); err != nil {
			log.Logger.Error(err.Error())
			payFailedChan <- msg
			continue
		}
		// 用户积分修改成功后，更新订单状态为已支付
		var shopClient shop.DaoShopClient
		shopClient, err = grpc.GetShopClient()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var reqOrder = &shop.ReqCommodityOrderStatusUpdateDao{
			AccountId:   msg.AccountId,
			Id:          msg.OrderId,
			OrderStatus: constant.Paid,
		}
		if _, err = shopClient.CommodityOrderStatusUpdateDao(context.Background(), reqOrder); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
	}
}

// 处理支付失败的问题
func handlePayFailed() {
	for {
		select {
		case msg := <-payFailedChan:
			shopClient, err := grpc.GetShopClient()
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			var reqOrder = &shop.ReqCommodityOrderStatusUpdateDao{
				AccountId:   msg.AccountId,
				Id:          msg.OrderId,
				OrderStatus: constant.PayFailed,
			}
			if _, err = shopClient.CommodityOrderStatusUpdateDao(context.Background(), reqOrder); err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		}
	}
}
