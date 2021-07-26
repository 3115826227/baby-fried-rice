package service

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

func InitBackend() {
	conf := config.GetConfig()
	if err := NewConsume(conf.NSQ.Topics.UserCoin, runUserCoinConsume); err != nil {
		return
	}
}

func NewConsume(consume config.TopicConsume, handle func(mq interfaces.MQ)) (err error) {
	conf := config.GetConfig()
	var mq = nsq.InitNSQMQ(conf.NSQ.Addr)
	err = mq.NewConsumer(consume.Topic, consume.Channel)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	go handle(mq)
	return
}

// 处理用户积分变动的消息
func runUserCoinConsume(mq interfaces.MQ) {
	for {
		value, err := mq.Consume()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var msg models.UserCoinChangeMQMessage
		if err = json.Unmarshal([]byte(value), &msg); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var userClient user.DaoUserClient
		userClient, err = grpc.GetUserClient()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var reqCoin = &user.ReqUserCoinLogAddDao{
			AccountId: msg.AccountId,
			Coin:      msg.Coin,
			CoinType:  msg.CoinType,
		}
		if _, err = userClient.UserCoinLogAddDao(context.Background(), reqCoin); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		// 用户积分修改成功后，更新订单状态
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
