package models

import "encoding/json"

type OrderCommodity struct {
	CommodityId string `json:"commodity_id"`
	PaymentType int64  `json:"payment_type"`
	PayedPrice  int64  `json:"payed_price"`
	PayedCoin   int64  `json:"payed_coin"`
}

func (oc *OrderCommodity) ToString() string {
	data, _ := json.Marshal(oc)
	return string(data)
}
