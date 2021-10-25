package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
	"baby-fried-rice/internal/pkg/module/shopDao/db"
	"baby-fried-rice/internal/pkg/module/shopDao/log"
	"baby-fried-rice/internal/pkg/module/shopDao/query"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type ShopService struct {
}

func CommodityModelToRpc(commodity tables.Commodity) *shop.CommodityQueryDao {
	return &shop.CommodityQueryDao{
		Id:       commodity.ID,
		Name:     commodity.Name,
		Title:    commodity.Title,
		Describe: commodity.Describe,
		SellType: int64(commodity.SellType),
		Price:    commodity.Price,
		Coin:     commodity.Coin,
		MainImg:  commodity.MainImg,
	}
}

func CommodityOrderModelToRpc(commodityOrder tables.CommodityOrder) *shop.CommodityOrderQueryDao {
	return &shop.CommodityOrderQueryDao{
		CommodityOrder: &shop.CommodityOrderBaseDao{
			Id:              commodityOrder.ID,
			AccountId:       commodityOrder.AccountId,
			PaymentType:     commodityOrder.PaymentType,
			TotalPrice:      commodityOrder.TotalPrice,
			TotalCoin:       commodityOrder.TotalCoin,
			Status:          int64(commodityOrder.Status),
			CreateTimestamp: commodityOrder.CreatedAt.Unix(),
			UpdateTimestamp: commodityOrder.UpdatedAt.Unix(),
		},
	}
}

func (service *ShopService) CommodityQueryDao(ctx context.Context, req *shop.ReqCommodityQueryDao) (resp *shop.RspCommodityQueryDao, err error) {
	var (
		commodities []tables.Commodity
		total       int64
	)
	if commodities, total, err = query.GetCommodities(req.Page, req.PageSize, ""); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*shop.CommodityQueryDao, 0)
	for _, c := range commodities {
		list = append(list, CommodityModelToRpc(c))
	}
	resp = &shop.RspCommodityQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

func (service *ShopService) CommodityDetailQueryDao(ctx context.Context, req *shop.ReqCommodityDetailQueryDao) (resp *shop.RspCommodityDetailQueryDao, err error) {
	var commodity tables.Commodity
	if commodity, err = query.GetCommodity(req.CommodityId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var relations []tables.CommodityImageRel
	if relations, err = query.GetCommodityImageRelation(req.CommodityId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var images = make([]string, 0)
	for _, rel := range relations {
		images = append(images, rel.Image)
	}
	resp = &shop.RspCommodityDetailQueryDao{
		Commodity: CommodityModelToRpc(commodity),
		Images:    images,
	}
	return
}

func (service *ShopService) CommodityCartUpdateDao(ctx context.Context, req *shop.ReqCommodityCartUpdateDao) (empty *emptypb.Empty, err error) {
	var updateMap = map[string]interface{}{
		"count":            req.UpdateCount,
		"update_timestamp": time.Now().Unix(),
	}
	// todo 考虑商品库存的问题
	if err = db.GetDB().GetDB().Model(&tables.CommodityCartRel{}).Where("account_id = ? and commodity_id = ?",
		req.AccountId, req.CommodityId).Updates(updateMap).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *ShopService) CommodityCartSelectDao(ctx context.Context, req *shop.ReqCommodityCartSelectDao) (empty *emptypb.Empty, err error) {
	var relations []tables.CommodityCartRel
	if relations, err = query.GetCommodityCartById(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var selectIds, unselectIds []string
	var selectIdMap = make(map[string]struct{})
	for _, id := range req.SelectedCommodityIds {
		selectIdMap[id] = struct{}{}
	}
	for _, rel := range relations {
		if _, exist := selectIdMap[rel.CommodityId]; !exist {
			if rel.Selected {
				unselectIds = append(unselectIds, rel.CommodityId)
			}
		} else {
			if !rel.Selected {
				selectIds = append(selectIds, rel.CommodityId)
			}
		}
	}
	var tx = db.GetDB().GetDB().Begin()
	var now = time.Now().Unix()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	if err = tx.Model(&tables.CommodityCartRel{}).Where("account_id = ? and commodity_id in (?)",
		req.AccountId, selectIds).Updates(map[string]interface{}{
		"selected":         true,
		"update_timestamp": now,
	}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = tx.Model(&tables.CommodityCartRel{}).Where("account_id = ? and commodity_id in (?)",
		req.AccountId, unselectIds).Updates(map[string]interface{}{
		"selected":         false,
		"update_timestamp": now,
	}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *ShopService) CommodityCartQueryDao(ctx context.Context, req *shop.ReqCommodityCartQueryDao) (resp *shop.RspCommodityCartQueryDao, err error) {
	var relations []tables.CommodityCartRel
	if relations, err = query.GetCommodityCartById(req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var ids = make([]string, 0)
	for _, rel := range relations {
		ids = append(ids, rel.CommodityId)
	}
	var commodities []tables.Commodity
	if commodities, err = query.GetCommoditiesByIds(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var commodityMap = make(map[string]tables.Commodity)
	for _, c := range commodities {
		commodityMap[c.ID] = c
	}
	var list = make([]*shop.CommodityCartDao, 0)
	for _, rel := range relations {
		var cart = &shop.CommodityCartDao{
			Commodity: CommodityModelToRpc(commodityMap[rel.CommodityId]),
			Count:     rel.Count,
			Selected:  rel.Selected,
		}
		list = append(list, cart)
	}
	resp = &shop.RspCommodityCartQueryDao{
		AccountId: req.AccountId,
		List:      list,
	}
	return
}

func (service *ShopService) CommodityCartDeleteDao(ctx context.Context, req *shop.ReqCommodityCartDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("account_id = ? and commodity_id = ?", req.AccountId,
		req.CommodityId).Delete(&tables.CommodityCartRel{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *ShopService) CommodityOrderAddDao(ctx context.Context, req *shop.ReqCommodityOrderAddDao) (empty *emptypb.Empty, err error) {
	var now = time.Now()
	var orderCommodities = make([]models.OrderCommodity, 0)
	for _, c := range req.OrderCommodities {
		var oc = models.OrderCommodity{
			CommodityId: c.CommodityId,
			PaymentType: c.PaymentType,
			PayedPrice:  c.PayedPrice,
			PayedCoin:   c.PayedCoin,
		}
		orderCommodities = append(orderCommodities, oc)
	}
	var data []byte
	if data, err = json.Marshal(orderCommodities); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var order = tables.CommodityOrder{
		AccountId:   req.AccountId,
		TotalPrice:  req.TotalPrice,
		TotalCoin:   req.TotalCoin,
		Commodities: string(data),
		Status:      constant.Submitted,
	}
	order.ID = handle.GenerateSerialNumberByLen(12)
	order.CreatedAt, order.UpdatedAt = now, now
	return
}

func (service *ShopService) CommodityOrderQueryDao(ctx context.Context, req *shop.ReqCommodityOrderQueryDao) (resp *shop.RspCommodityOrderQueryDao, err error) {
	var (
		commodityOrders []tables.CommodityOrder
		total           int64
	)
	if commodityOrders, total, err = query.GetCommodityOrders(req.Page, req.PageSize, req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*shop.CommodityOrderQueryDao, 0)
	for _, co := range commodityOrders {
		var commodityOrder = CommodityOrderModelToRpc(co)
		var orderCommodities []models.OrderCommodity
		if err = json.Unmarshal([]byte(co.Commodities), &orderCommodities); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var ids = make([]string, 0)
		var ocMap = make(map[string]models.OrderCommodity)
		for _, c := range orderCommodities {
			ocMap[c.CommodityId] = c
			ids = append(ids, c.CommodityId)
		}
		var commodities []tables.Commodity
		if commodities, err = query.GetCommoditiesByIds(ids); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var commoditiesDao = make([]*shop.OrderCommodityDao, 0)
		for _, c := range commodities {
			var commodityDao = new(shop.OrderCommodityDao)
			commodityDao.Commodity = CommodityModelToRpc(c)
			commoditiesDao = append(commoditiesDao, commodityDao)
		}
		commodityOrder.Commodities = commoditiesDao
		list = append(list, commodityOrder)
	}
	resp = &shop.RspCommodityOrderQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

func (service *ShopService) CommodityOrderDetailQueryDao(ctx context.Context, req *shop.ReqCommodityOrderDetailQueryDao) (resp *shop.RspCommodityOrderDetailQueryDao, err error) {
	// 查询指定订单
	var co tables.CommodityOrder
	if co, err = query.GetCommodityOrderById(req.Id, req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 查询订单下的商品列表
	var orderCommodities []models.OrderCommodity
	if err = json.Unmarshal([]byte(co.Commodities), &orderCommodities); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var ids = make([]string, 0)
	var ocMap = make(map[string]models.OrderCommodity)
	for _, c := range orderCommodities {
		ocMap[c.CommodityId] = c
		ids = append(ids, c.CommodityId)
	}
	// 查询商品信息
	var commodities []tables.Commodity
	if commodities, err = query.GetCommoditiesByIds(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var orderCommoditiesDao = make([]*shop.OrderCommodityDao, 0)
	for _, c := range commodities {
		var orderCommodityDao = &shop.OrderCommodityDao{
			Commodity:   CommodityModelToRpc(c),
			PaymentType: ocMap[c.ID].PaymentType,
			PayedPrice:  ocMap[c.ID].PayedPrice,
			PayedCoin:   ocMap[c.ID].PayedCoin,
		}
		orderCommoditiesDao = append(orderCommoditiesDao, orderCommodityDao)
	}
	resp = &shop.RspCommodityOrderDetailQueryDao{
		CommodityOrder:   CommodityOrderModelToRpc(co).CommodityOrder,
		CommodityDetails: orderCommoditiesDao,
	}
	return
}

func (service *ShopService) CommodityOrderStatusUpdateDao(ctx context.Context, req *shop.ReqCommodityOrderStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(&tables.CommodityOrder{}).Where("id = ? and account_id = ?", req.Id, req.AccountId).Update("status", req.OrderStatus).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *ShopService) CommodityOrderDeleteDao(ctx context.Context, req *shop.ReqCommodityOrderDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("id = ? and account_id = ?", req.Id, req.AccountId).Delete(&tables.CommodityOrder{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}
