package constant

type SellType int64

const (
	// 金额支付
	SellForMoney SellType = 1
	// 积分兑换
	SellForCoin = 2
	// 参与积分抽奖
	SellForCoinLottery = 3
)

type OrderStatus int64

const (
	Submitted OrderStatus = 1 // 已提交待支付
	Paying                = 2 // 支付中请稍等
	Paid                  = 3 // 支付成功待发货
	Shipping              = 4 // 发货中
	Shipped               = 5 // 已发货待接收
	Completed             = 6 // 已完成

	CancellingReview = 101 // 取消审核中
	Cancelled        = 102 // 已取消

	PayFailed        = 201 // 支付失败待重新支付
	WaitPayedTimeout = 202 // 等待支付超时（订单创建30分钟之内未支付，将关闭订单）
)
