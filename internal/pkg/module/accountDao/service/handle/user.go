package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UsersHandle(c *gin.Context) {
	pageInfo, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		offset = (pageInfo.Page - 1) * pageInfo.PageSize
		limit  = pageInfo.PageSize
		users  = make([]tables.AccountUserDetail, 0)
		total  int64
	)
	if err = db.GetDB().GetDB().Model(&tables.AccountUserDetail{}).Count(&total).Order("update_time").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	list := make([]interface{}, 0)
	for _, user := range users {
		list = append(list, user)
	}
	handle.SuccessListResp(c, "", list, total, pageInfo)
}
