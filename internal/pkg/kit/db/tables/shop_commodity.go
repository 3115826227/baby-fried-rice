package tables

import (
	"baby-fried-rice/internal/pkg/kit/constant"
)

// 商品基础表
type Commodity struct {
	CommonField
	// 商品名称
	Name string `gorm:"column:name"`
	// 商品标题
	Title string `gorm:"column:title"`
	// 商品描述
	Describe string `gorm:"column:describe"`
	// 售卖方式
	SellType constant.SellType `gorm:"column:sell_type"`
	// 商品价格
	Price int64 `gorm:"column:price"`
	// 积分兑换数量
	Coin int64 `gorm:"column:coin"`
	// 商品主图片
	MainImg string `gorm:"column:main_img"`
	// 商品状态 1-未上架 2-上架
	Status int64 `gorm:"column:status"`
}

func (table *Commodity) TableName() string {
	return "baby_shop_commodity"
}

// 商品图片关系表
type CommodityImageRel struct {
	CommodityID     string `gorm:"commodity_id"`
	Image           string `gorm:"image"`
	CreateTimestamp int64  `gorm:"create_timestamp"`
}

func (table *CommodityImageRel) TableName() string {
	return "baby_shop_commodity_image_relation"
}
