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

// 用户信息列表查询
func UserHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		details []tables.AccountUserDetail
		total   int64
	)
	details, total, err = query.GetUsers(reqPage.Page, reqPage.PageSize)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]rsp.UserBackendResp, 0)
	for _, d := range details {
		var detail = rsp.UserBackendResp{
			AccountId:    d.AccountID,
			HeadImgUrl:   d.HeadImgUrl,
			Username:     d.Username,
			SchoolId:     d.SchoolId,
			Gender:       d.Gender,
			Age:          d.Age,
			Phone:        d.Phone,
			RegisterTime: d.CreatedAt.Unix(),
		}
		list = append(list, detail)
	}
	var response = rsp.UserBackendListResp{
		List:     list,
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
		Total:    total,
	}
	handle.SuccessResp(c, "", response)
}
