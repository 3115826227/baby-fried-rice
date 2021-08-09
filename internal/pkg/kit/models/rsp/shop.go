package rsp

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
)

func CommodityModelToRsp(commodity tables.Commodity) Commodity {
	return Commodity{
		Id:         commodity.ID,
		Name:       commodity.Name,
		Title:      commodity.Title,
		Describe:   commodity.Describe,
		SellType:   commodity.SellType,
		Price:      commodity.Price,
		Coin:       commodity.Coin,
		MainImg:    commodity.MainImg,
		CreateTime: commodity.CreatedAt.Unix(),
		UpdateTime: commodity.UpdatedAt.Unix(),
	}
}

func CommodityRpcToRsp(commodity *shop.CommodityQueryDao) Commodity {
	return Commodity{
		Id:       commodity.Id,
		Name:     commodity.Name,
		Title:    commodity.Title,
		Describe: commodity.Describe,
		SellType: constant.SellType(commodity.SellType),
		Price:    commodity.Price,
		Coin:     commodity.Coin,
		MainImg:  commodity.MainImg,
	}
}

type Commodity struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Title      string            `json:"title"`
	Describe   string            `json:"describe"`
	SellType   constant.SellType `json:"sell_type"`
	Price      int64             `json:"price"`
	Coin       int64             `json:"coin"`
	MainImg    string            `json:"main_img"`
	CreateTime int64             `json:"create_time,omitempty"`
	UpdateTime int64             `json:"update_time,omitempty"`
}

type CommodityDetailResp struct {
	Commodity
	Images []string `json:"images"`
}

type CommodityOrdersResp struct {
	List     []CommodityOrder `json:"list"`
	Page     int64            `json:"page"`
	PageSize int64            `json:"page_size"`
	Total    int64            `json:"total"`
}

type CommodityOrder struct {
	CommodityOrderBase
	Commodities []Commodity `json:"commodities"`
}

type CommodityOrderBase struct {
	Id              string `json:"id"`
	AccountId       string `json:"account_id"`
	PaymentType     int64  `json:"payment_type"`
	TotalPrice      int64  `json:"total_price"`
	TotalCoin       int64  `json:"total_coin"`
	Status          int64  `json:"status"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}

type CommodityOrderDetailResp struct {
	CommodityOrderBase
	Commodities []OrderCommodityDetail `json:"commodities"`
}

type OrderCommodityDetail struct {
	Commodity
	PaymentType int64 `json:"payment_type"`
	PayedPrice  int64 `json:"payed_price"`
	PayedCoin   int64 `json:"payed_coin"`
}

type PayOrderCommodityResp struct {
	// 操作结果
	Ok bool `json:"ok"`
	// 描述信息
	Describe string `json:"describe"`
	// 操作时间
	Timestamp int64 `json:"timestamp"`
}
