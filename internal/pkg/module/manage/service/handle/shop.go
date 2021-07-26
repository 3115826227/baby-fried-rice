package handle

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 添加商品
func AddCommodityHandle(c *gin.Context) {
	var req requests.ReqAddCommodity
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var now = time.Now()
	var commodity = tables.Commodity{
		Name:     req.Name,
		Title:    req.Title,
		Describe: req.Describe,
		SellType: req.SellType,
		Price:    req.Price,
		Coin:     req.Coin,
		MainImg:  req.MainImg,
	}
	commodity.ID = handle.GenerateSerialNumberByLen(10)
	commodity.CreatedAt, commodity.UpdatedAt = now, now
	if err := db.GetShopDB().CreateObject(&commodity); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查询商品列表
func CommodityHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		commodities []tables.Commodity
		total       int64
	)
	var param = query.CommoditiesQueryParam{
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
	}
	commodities, total, err = query.GetCommodities(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]rsp.Commodity, 0)
	for _, commodity := range commodities {
		list = append(list, rsp.CommodityModelToRsp(commodity))
	}
	var response = rsp.CommoditiesResp{
		List:     list,
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
		Total:    total,
	}
	handle.SuccessResp(c, "", response)
}
