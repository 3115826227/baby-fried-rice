package tables

import "baby-fried-rice/internal/pkg/kit/constant"

// 商品订单基础表
type CommodityOrder struct {
	CommonField

	// 购买者用户id
	AccountId string `gorm:"account_id"`
	// 支付方式
	PaymentType int64 `gorm:"payment_type"`
	// 总支付金额
	TotalPrice int64 `gorm:"total_price"`
	// 总消耗积分
	TotalCoin int64 `gorm:"total_coin"`
	// 订单商品列表信息
	Commodities string `gorm:"commodities"`
	// 订单状态
	Status constant.OrderStatus `gorm:"status"`
}

func (table *CommodityOrder) TableName() string {
	return "baby_shop_commodity_order"
}
