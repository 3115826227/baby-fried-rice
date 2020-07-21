package model

import (
	"shop/src/log"
	"shop/src/model/db"
	"time"
)

func init() {
	err := db.GetDB().AutoMigrate().Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type IntCommonField struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"updated_at"`
}

type BigIntCommonField struct {
	ID        int64     `gorm:"AUTO_INCREMENT;column:id;" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"updated_at"`
}

type StringCommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"updated_at"`
}

//平台用户信息
type UserMetaPlatform struct {
	UserId   string
	Nickname string

	Coins       string //积分
	Credibility int    //购买信誉
}

//用户收货地址
type UserShopReceiveAddress struct {
	UserId            string
	ReceiveUserName   string //收货人姓名
	ReceiveAddressSeq int    //收货地址序列号
	Province          string
	City              string
	Local             string
	ReceiveUserPhone  string //收货人联系方式
}

//用户购物车
type UserShopCart struct {
	UserId         string
	CommoditySpuId string //商品spu_id
	Number         int    //该商品数量
}

//用户会员
type UserShopMember struct {
	IntCommonField
	UserId             string //用户user_id
	EffectiveTimestamp int64  //会员有效时间
	FailureTimestamp   int64  //会员失效时间
	ExpireTimestamp    int64  //会员有效持续时间
}

//平台优惠券模板
type PlatformCouponTemplate struct {
	IntCommonField
	Name            string //模板名称
	CouponType      int    //优惠券类型：1-满减券，2-折扣券
	ReceiveRule     int    //领取规则
	UserReceiveRule int    //使用规则
}

//平台优惠券发放信息
type PlatformCouponInfo struct {
}

//平台优惠券用户领取信息
type PlatformCouponInfoUser struct {
	PlatformCouponInfoId int
	UserId               string
	UseStatus            int //优惠券的使用状态: 未使用/已使用/已过期
}

//平台满减劵
type PlatformFullReduceCoupon struct {
	PlatformCouponTemplateId int
	FullReduce               bool //true-满减劵，false-立减劵，无门槛
	FullMoney                int  //满金额
	ReduceMoney              int  //减金额
}

//平台折扣券
type PlatformShopDiscountCoupon struct {
	PlatformCouponTemplateId int
	Discount                 int //折扣，比如打八折
	ConditionType            int //折扣使用条件，0-无条件，1-金额条件，2-商品数量条件
	MoneyCondition           int //折扣使用金额条件，比如满200 打八折
	CountCondition           int //折扣使用数量条件，比如满两件 打八折
}

//店铺优惠券模板
type ShopCouponTemplate struct {
	IntCommonField
	ShopId     int    //店铺id
	Name       string //模板名称
	CouponType int    //优惠券类型：1-满减券，2-折扣券
}

//店铺满减劵
type ShopFullReduceCoupon struct {
	ShopCouponTemplateId int
	FullReduce           bool //true-满减劵，false-立减劵，无门槛
	FullMoney            int  //满金额
	ReduceMoney          int  //减金额
}

//店铺折扣券
type ShopDiscountCoupon struct {
	ShopCouponTemplateId int
	Discount             int //折扣，比如打八折
	ConditionType        int //折扣使用条件，0-无条件，1-金额条件，2-商品数量条件
	MoneyCondition       int //折扣使用金额条件，比如满200 打八折
	CountCondition       int //折扣使用数量条件，比如满两件 打八折
}

//店铺优惠券发放信息
type ShopCouponInfo struct {
	IntCommonField
	Name                  string //优惠券名称
	ShopCouponTemplateId  int    //优惠券模板id
	Count                 int64  //优惠券发放数量
	LimitUserCount        int    //用户限领张数
	SupplyStartTimestamp  int64  //优惠券发放开始时间
	SupplyEndTimestamp    int64  //优惠券发放截止时间
	UseEffectiveTimestamp int64  //优惠券生效时间
	UseFailureTimestamp   int64  //优惠券失效时间
}

//店铺优惠券用户领取信息
type ShopCouponInfoUser struct {
	ShopCouponInfoId int
	UserId           string
	UseStatus        int //优惠券的使用状态: 未使用/已使用/已过期
}

//商品Spu
type CommodityMeta struct {
	BigIntCommonField

	CategoryId      int   //分类
	BrandId         int   //品牌
	ShopId          int   //商家
	ListStatus      bool  //是否上架
	CommodityUnitId int   //商品单元
	Sales           int64 //销售量

	Describe   string //商品描述
	MetaDetail string //商品属性描述
}

//商品Sku
type CommoditySku struct {
	StringCommonField
	CommoditySpuId    int64  //商品spu_id
	StockStatus       bool   //商品是否有库存
	CommoditySpecs    string //商品sku属性描述
	CommoditySpecsSeq int    //商品sku属性序列号
}

//商品价格
type CommodityPrice struct {
	CommoditySkuId  string
	PriceCategoryId int   //商品价格类别id
	CurrentUse      bool  //当前是否使用该价格
	Price           int64 //商品价格 单位：分
}

/*
	商品价格类别
	  * 原价
      * 活动价
      * 秒杀价
*/
type CommodityPriceCategory struct {
	IntCommonField
	Category string //商品价格类别名
}

//商品库存
type CommodityStock struct {
	CommoditySkuId  string
	CommodityUnitID int
	Stock           int64
}

//商品运费计算方式
type CommodityFreight struct {
	IntCommonField
	Name      string //商品运费计算方式名称：免运费/49包邮等
	FreePrice int64  //免邮价格  单元：分
	Freight   int64  //邮费价格  单元：分
}

//商品类别
type CommodityCategory struct {
	IntCommonField
	Category string //商品类别名
	ParentId int    //商品父类id
}

//商品品牌
type CommodityBrand struct {
	IntCommonField
	Brand string //商品品牌名
}

//商品元数据属性
type CommodityMetaAttribute struct {
	IntCommonField
	CommoditySpuId     int64  //商品spu_id
	AttributeKey       string //商品元数据属性名
	AttributeValueType int    //商品元数据属性值类型
	AttributeValue     string //商品元数据属性值
}

//商品销售属性Key
type CommodityAttributeKey struct {
	IntCommonField
	CommoditySpuId int64  //商品spu_id
	AttributeKey   string //商品销售属性key
}

//商品销售属性Value
type CommodityAttributesValue struct {
	BigIntCommonField
	AttributeKeyId int    //商品销售属性Key id
	AttributeValue string //商品销售属性Value
}

//商品单位
type CommodityUnit struct {
	IntCommonField
	Unit string //商品单位名 例如：件/套/双/瓶/箱等
}

//商品标签
type CommodityLabel struct {
	IntCommonField
	Label string //商品标签名
}

//商品评价
type CommodityEvaluation struct {
	CommoditySkuId    int64
	EvaluationLevelId int
	Content           string
}

//商品评价等级
type CommodityEvaluationLevel struct {
	IntCommonField
	Level string //商品评价等级名 例如：五星好评
}

//商品提问
type CommodityQuestion struct {
	IntCommonField
	CommoditySpuId int64
	Question       string
}

//商品回答
type CommodityAnswer struct {
	IntCommonField
	QuestionId int
	Answer     string
}

//商品Sku收藏
type CommodityCollect struct {
	IntCommonField
	UserId         string
	Username       string
	CommoditySkuId string //商品sku_id
}

//商品Spu关注/举报
type CommodityAttentionReport struct {
	IntCommonField
	UserId         string
	Username       string
	Attention      bool  //商品spu关注
	Report         bool  //商品spu举报
	CommoditySpuId int64 //商品spu_id
}

type CommodityReportType struct {
	IntCommonField
	Reason string
}

type ShopLogin struct {
	IntCommonField
	LoginName string `gorm:"column:login_name;unique;not null" json:"login_name"`
	Password  string `gorm:"column:password;not null"`
}

//店铺
type Shop struct {
	ShopId   int    `gorm:"column:shop_id;unique;not null" json:"shop_id"`
	Name     string `gorm:"column:name;unique;not null" json:"name"`
	Describe string `gorm:"column:describe;" json:"describe"`
	Phone    string `gorm:"column:phone;unique;not null" json:"phone"`
	Status   bool   `gorm:"column:status;" json:"status"`
}

type ShopFaction struct {
	ShopId           int
	ServiceFaction   int //店铺服务评分
	CommodityFaction int //店铺商品评分
	LogisticsFaction int //店铺物流评分
}

//店铺经营商品类别
type ShopManageCategory struct {
	CommodityCategoryId int `gorm:"column:commodity_category_id;unique_index:shop_manage_commodity_category_shop" json:"commodity_category_id"`
	ShopId              int `gorm:"column:shop_id;unique_index:shop_manage_commodity_category_shop" json:"shop_id"`
}

//店铺经营品牌
type ShopManageBrand struct {
	CommodityBrandId int `gorm:"column:commodity_brand_id;unique_index:shop_manage_commodity_brand_shop" json:"brand_id"`
	ShopId           int `gorm:"column:shop_id;unique_index:shop_manage_commodity_brand_shop" json:"shop_id"`
}

//店铺关注/举报
type ShopAttentionReport struct {
	IntCommonField
	UserId    string
	Username  string
	Attention bool //商品spu关注
	Report    bool //商品spu举报
	ShopId    int  //商家id
}

//店铺举报理由
type ShopAttentionReportType struct {
	IntCommonField
	Reason string
}

//订单
type Order struct {
}

//订单评价
type OrderEvaluation struct {
}

//订单评价等级
type OrderEvaluationLevel struct {
}

//订单支付方式
type OrderPay struct {
}

//物流
type Logistics struct {
}
