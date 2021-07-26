package requests

import "baby-fried-rice/internal/pkg/kit/constant"

// 商品添加请求
type ReqAddCommodity struct {
	// 商品名称
	Name string `json:"name"`
	// 商品标题
	Title string `json:"title"`
	// 商品描述
	Describe string `json:"describe"`
	// 售卖方式
	SellType constant.SellType `json:"sell_type"`
	// 商品价格
	Price int64 `json:"price"`
	// 积分兑换数量
	Coin int64 `json:"coin"`
	// 商品主图片
	MainImg string `json:"main_img"`
}

// 下单请求
type ReqAddCommodityOrder struct {
	Commodities []OrderCommodity `json:"commodities" binding:"required"`
	TotalPrice  int64            `json:"total_price"`
	TotalCoin   int64            `json:"total_coin"`
}

type OrderCommodity struct {
	CommodityId string `json:"commodity_id" binding:"required"`
	PaymentType int64  `json:"payment_type"`
	PayedPrice  int64  `json:"payed_price"`
	PayedCoin   int64  `json:"payed_coin"`
}

// 支付订单
type ReqPayCommodityOrder struct {
	ID          string `json:"id" binding:"required"`
	PaymentType int64  `json:"payment_type" binding:"required"`
}
