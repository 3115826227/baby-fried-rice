package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/config"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 更新商品
func UpdateCommodityHandle(c *gin.Context) {
	var req requests.ReqUpdateCommodity
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var commodity = tables.Commodity{
		Name:     req.Name,
		Title:    req.Title,
		Describe: req.Describe,
		SellType: req.SellType,
		Price:    req.Price,
		Coin:     req.Coin,
		MainImg:  req.MainImg,
	}
	if req.Status != nil {
		commodity.Status = *req.Status
	}
	commodity.ID = req.Id
	commodity.UpdatedAt = time.Now()
	if err := db.GetShopDB().GetDB().Updates(&commodity).Error; err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查询商品列表
func CommodityHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var (
		commodities []tables.Commodity
		total       int64
	)
	var param = query.CommoditiesQueryParam{
		SellType: c.Query("sell_type"),
		LikeName: c.Query(handle.QueryLikeName),
		Status:   c.Query("status"),
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
	}
	commodities, total, err = query.GetCommodities(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, commodity := range commodities {
		list = append(list, rsp.CommodityModelToRsp(commodity))
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}

// 删除商品
func DeleteCommodityHandle(c *gin.Context) {
	id := c.Query("id")
	var commodity tables.Commodity
	if err := db.GetShopDB().GetDB().Where("id = ?", id).First(&commodity).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var relations []tables.CommodityImageRel
	if err := db.GetShopDB().GetDB().Where("commodity_id = ?", id).Find(&relations).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var imgs = make([]string, 0)
	imgs = append(imgs, commodity.MainImg)
	for _, rel := range relations {
		imgs = append(imgs, rel.Image)
	}
	var err error
	tx := db.GetShopDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit().Error; err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	if err = tx.Delete(&commodity).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	if err = tx.Delete(&tables.CommodityImageRel{}, "commodity_id = ?", id).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	go func() {
		// 删除云存储中的图片地址
		var topic = config.GetConfig().MessageQueue.PublishTopics.DeleteFile
		for _, img := range imgs {
			var info = models.DeleteFileMessageQueueInfo{
				FileValueType: models.FileDownUrl,
				FileValue:     img,
			}
			if err = mq.Send(topic, info.ToString()); err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		}
	}()
	handle.SuccessResp(c, "", nil)
}
