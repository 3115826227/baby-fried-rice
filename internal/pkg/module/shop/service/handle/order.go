package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/shop/config"
	"baby-fried-rice/internal/pkg/module/shop/grpc"
	"baby-fried-rice/internal/pkg/module/shop/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CommodityOrderBaseRpcToRsp(commodityOrderBase *shop.CommodityOrderBaseDao) rsp.CommodityOrderBase {
	return rsp.CommodityOrderBase{
		Id:              commodityOrderBase.Id,
		AccountId:       commodityOrderBase.AccountId,
		PaymentType:     commodityOrderBase.PaymentType,
		TotalPrice:      commodityOrderBase.TotalPrice,
		TotalCoin:       commodityOrderBase.TotalCoin,
		Status:          commodityOrderBase.Status,
		CreateTimestamp: commodityOrderBase.CreateTimestamp,
		UpdateTimestamp: commodityOrderBase.UpdateTimestamp,
	}
}

// 商品订单添加
func CommodityOrderAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqAddCommodityOrder
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var orderCommodities = make([]*shop.OrderCommodityAddDao, 0)
	for _, commodity := range req.Commodities {
		orderCommodities = append(orderCommodities, &shop.OrderCommodityAddDao{
			CommodityId: commodity.CommodityId,
			PaymentType: commodity.PaymentType,
			PayedPrice:  commodity.PayedPrice,
			PayedCoin:   commodity.PayedCoin,
		})
	}
	var reqShop = &shop.ReqCommodityOrderAddDao{
		AccountId:        userMeta.AccountId,
		OrderCommodities: orderCommodities,
		TotalPrice:       req.TotalPrice,
		TotalCoin:        req.TotalCoin,
	}
	_, err = shopClient.CommodityOrderAddDao(context.Background(), reqShop)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 支付订单
func CommodityOrderPayHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqPayCommodityOrder
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqShop = &shop.ReqCommodityOrderDetailQueryDao{
		Id:        req.ID,
		AccountId: userMeta.AccountId,
	}
	var resp *shop.RspCommodityOrderDetailQueryDao
	resp, err = shopClient.CommodityOrderDetailQueryDao(context.Background(), reqShop)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var response = rsp.PayOrderCommodityResp{
		Timestamp: time.Now().UnixNano(),
	}
	switch constant.OrderStatus(resp.CommodityOrder.Status) {
	case constant.Submitted:
	case constant.PayFailed:
	case constant.Paying:
		response.Describe = "订单支付请求已经提交，正在处理，请稍等片刻"
		handle.SuccessResp(c, "", response)
		return
	case constant.Cancelled:
		response.Describe = "订单已经提交取消请求或取消请求正在审核中，无法进行支付"
		handle.SuccessResp(c, "", response)
		return
	case constant.CancellingReview:
		response.Describe = "订单已经提交取消请求或取消请求正在审核中，无法进行支付"
		handle.SuccessResp(c, "", response)
		return
	case constant.WaitPayedTimeout:
		response.Describe = "订单等待支付时间过长已关闭，无法进行支付"
		handle.SuccessResp(c, "", response)
		return
	default:
		response.Describe = "订单已支付成功，请勿重新提交"
		handle.SuccessResp(c, "", response)
		return
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var respCoin *user.RspUserCoinDao
	respCoin, err = userClient.UserCoinDao(context.Background(), &user.ReqUserCoinDao{AccountId: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if respCoin.Coin < resp.CommodityOrder.TotalCoin {
		response.Describe = "积分不足，无法支付"
		handle.SuccessResp(c, "", response)
		return
	}
	var reqShopStatus = &shop.ReqCommodityOrderStatusUpdateDao{
		AccountId:   userMeta.AccountId,
		Id:          req.ID,
		OrderStatus: constant.Paying,
	}
	_, err = shopClient.CommodityOrderStatusUpdateDao(context.Background(), reqShopStatus)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 如果是积分兑换，则发送积分变动消息给消息队列，交给下游处理
	if resp.CommodityOrder.TotalCoin != 0 {
		var mqMessage = models.UserCoinChangeMQMessage{
			AccountId: userMeta.AccountId,
			Coin:      -resp.CommodityOrder.TotalCoin,
			CoinType:  constant.ConsumeCoinType,
			OrderId:   resp.CommodityOrder.Id,
		}
		go func() {
			if err = mq.Send(config.GetConfig().MessageQueue.PublishTopics.UserCoin, mqMessage.ToString()); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}()
	}
	handle.SuccessResp(c, "", response)
}

// 商品订单列表查询
func CommodityOrderQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var shopClient shop.DaoShopClient
	shopClient, err = grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &shop.ReqCommodityOrderQueryDao{
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	var resp *shop.RspCommodityOrderQueryDao
	resp, err = shopClient.CommodityOrderQueryDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, commodityOrder := range resp.List {
		var commodities = make([]rsp.Commodity, 0)
		for _, commodity := range commodityOrder.Commodities {
			commodities = append(commodities, rsp.CommodityRpcToRsp(commodity.Commodity))
		}
		var co = rsp.CommodityOrder{
			CommodityOrderBase: CommodityOrderBaseRpcToRsp(commodityOrder.CommodityOrder),
			Commodities:        commodities,
		}
		list = append(list, co)
	}
	handle.SuccessListResp(c, "", list, resp.Total, reqPage.Page, reqPage.PageSize)
}

// 商品订单详细信息查询
func CommodityOrderDetailQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	id := c.Query("id")
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &shop.ReqCommodityOrderDetailQueryDao{
		Id:        id,
		AccountId: userMeta.AccountId,
	}
	var resp *shop.RspCommodityOrderDetailQueryDao
	resp, err = shopClient.CommodityOrderDetailQueryDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var orderCommodityDetails = make([]rsp.OrderCommodityDetail, 0)
	for _, commodityDetail := range resp.CommodityDetails {
		var orderCommodityDetail = rsp.OrderCommodityDetail{
			Commodity:   rsp.CommodityRpcToRsp(commodityDetail.Commodity),
			PaymentType: commodityDetail.PaymentType,
			PayedPrice:  commodityDetail.PayedPrice,
			PayedCoin:   commodityDetail.PayedCoin,
		}
		orderCommodityDetails = append(orderCommodityDetails, orderCommodityDetail)
	}
	var response = rsp.CommodityOrderDetailResp{
		CommodityOrderBase: CommodityOrderBaseRpcToRsp(resp.CommodityOrder),
		Commodities:        orderCommodityDetails,
	}
	handle.SuccessResp(c, "", response)
}

// 商品订单删除
func CommodityOrderDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	id := c.Query("id")
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var req = &shop.ReqCommodityOrderDeleteDao{
		AccountId: userMeta.AccountId,
		Id:        id,
	}
	_, err = shopClient.CommodityOrderDeleteDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
