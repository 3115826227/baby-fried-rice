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

// 空间列表查询
func SpaceHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		spaces []tables.Space
		total  int64
	)
	var param = query.SpacesQueryParam{
		Id:          c.Query(handle.QueryId),
		AccountId:   c.Query(handle.QueryAccountId),
		VisitorType: c.Query("visitor_type"),
		AuditStatus: c.Query("audit_status"),
		Page:        reqPage.Page,
		PageSize:    reqPage.PageSize,
	}
	spaces, total, err = query.GetSpaces(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	var idMap = make(map[string]struct{})
	for _, s := range spaces {
		idMap[s.Origin] = struct{}{}
	}
	var ids = make([]string, 0)
	for origin := range idMap {
		ids = append(ids, origin)
	}
	var idsMap = make(map[string]tables.AccountUserDetail)
	if idsMap, err = query.GetUsersByIds(ids); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	for _, s := range spaces {
		var user = idsMap[s.Origin]
		var space = rsp.AdminSpaceResp{
			Id:          s.ID,
			AuditStatus: s.AuditStatus,
			VisitorType: s.VisitorType,
			User: rsp.User{
				AccountID:  user.AccountID,
				Username:   user.Username,
				HeadImgUrl: user.HeadImgUrl,
				IsOfficial: user.IsOfficial,
			},
			CreateTime: s.CreatedAt.Unix(),
			UpdateTime: s.UpdatedAt.Unix(),
		}
		list = append(list, space)
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}

// 空间列表审核
func UpdateSpaceAuditHandle(c *gin.Context) {
	var req requests.ReqUpdateSpaceAudit
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var space tables.Space
	if err := db.GetSpaceDB().GetObject(map[string]interface{}{"id": req.SpaceId}, &space); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	if req.Audit != nil {
		space.AuditStatus = *req.Audit
		space.UpdatedAt = time.Now()
		if err := db.GetSpaceDB().UpdateObject(&space); err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	}
	handle.SuccessResp(c, "", nil)
}
