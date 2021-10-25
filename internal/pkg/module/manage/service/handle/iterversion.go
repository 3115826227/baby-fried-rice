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

// 迭代版本添加
func AddIterativeVersionHandle(c *gin.Context) {
	var req requests.ReqAddIterativeVersion
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var now = time.Now().Unix()
	var iv = tables.IterativeVersion{
		Version:         req.Version,
		Content:         req.Content,
		CreateTimestamp: now,
		UpdateTimestamp: now,
	}
	if err := db.GetAccountDB().CreateObject(&iv); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 迭代版本内容更新
func UpdateIterativeVersionHandle(c *gin.Context) {
	var req requests.ReqUpdateIterativeVersion
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	iv, err := query.GetIterativeVersionByVersion(req.Version)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if req.Content != nil {
		iv.Content = *req.Content
	}
	if req.Status != nil {
		iv.Status = *req.Status
	}
	iv.UpdateTimestamp = time.Now().Unix()
	if err = db.GetAccountDB().UpdateObject(&iv); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 迭代版本查询
func QueryIterativeVersionHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		ivs   []tables.IterativeVersion
		total int64
	)
	var param = query.IterativeVersionQueryParam{
		LikeVersion: c.Query("version"),
		Page:        reqPage.Page,
		PageSize:    reqPage.PageSize,
	}
	ivs, total, err = query.GetIterativeVersion(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, iv := range ivs {
		list = append(list, rsp.ManageIterativeVersion{
			BaseIterativeVersion: rsp.BaseIterativeVersion{
				Version:         iv.Version,
				Content:         iv.Content,
				UpdateTimestamp: iv.UpdateTimestamp,
			},
			CreateTimestamp: iv.CreateTimestamp,
			Status:          iv.Status,
		})
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}
