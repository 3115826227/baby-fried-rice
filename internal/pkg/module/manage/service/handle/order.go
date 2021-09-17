package handle

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 查询订单列表
func OrderHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		orders []tables.CommodityOrder
		total  int64
	)
	var param = query.OrderQueryParam{
		OrderId:   c.Query("order_id"),
		AccountId: c.Query(handle.QueryAccountId),
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	orders, total, err = query.GetOrders(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, order := range orders {
		list = append(list, rsp.CommodityOrderModelToRsp(order))
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}
