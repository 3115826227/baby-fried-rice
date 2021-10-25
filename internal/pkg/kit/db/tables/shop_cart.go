package tables

// 用户购物车表
type CommodityCartRel struct {
	// 用户id
	AccountId string `gorm:"column:account_id;unique_index:cart_account_commodity"`
	// 商品id
	CommodityId string `gorm:"column:commodity_id;unique_index:cart_account_commodity"`
	// 数量
	Count int64 `gorm:"column:count"`
	// 选择状态 0-未勾选 1-已勾选
	Selected bool `gorm:"column:selected"`
	// 更新时间
	UpdateTimestamp int64 `gorm:"column:update_timestamp"`
}

func (table *CommodityCartRel) TableName() string {
	return "baby_shop_commodity_cart_relation"
}
