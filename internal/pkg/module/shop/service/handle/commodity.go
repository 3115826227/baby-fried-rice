package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
	"baby-fried-rice/internal/pkg/module/shop/grpc"
	"baby-fried-rice/internal/pkg/module/shop/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CommodityQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = &shop.ReqCommodityQueryDao{
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	var resp *shop.RspCommodityQueryDao
	resp, err = shopClient.CommodityQueryDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, commodity := range resp.List {
		list = append(list, rsp.CommodityRpcToRsp(commodity))
	}
	handle.SuccessListResp(c, "", list, resp.Total, resp.Page, resp.PageSize)
}

func CommodityDetailQueryHandle(c *gin.Context) {
	id := c.Query("id")
	shopClient, err := grpc.GetShopClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *shop.RspCommodityDetailQueryDao
	resp, err = shopClient.CommodityDetailQueryDao(context.Background(),
		&shop.ReqCommodityDetailQueryDao{CommodityId: id})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var response = rsp.CommodityDetailResp{
		Commodity: rsp.CommodityRpcToRsp(resp.Commodity),
		Images:    resp.Images,
	}
	handle.SuccessResp(c, "", response)
}
