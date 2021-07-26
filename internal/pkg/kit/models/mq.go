package models

import "encoding/json"

// 用户积分变动信息
type UserCoinChangeMQMessage struct {
	// 变动用户
	AccountId string `json:"account_id"`
	// 变动积分
	Coin int64 `json:"coin"`
	// 积分使用类型
	CoinType int64 `json:"coin_type"`
	// 关联订单id
	OrderId string `json:"order_id"`
}

func (message *UserCoinChangeMQMessage) ToString() string {
	data, _ := json.Marshal(message)
	return string(data)
}

// 订单状态变动消息
type OrderStatusChangeMQMessage struct {
}
