package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 用户添加
func AddUserHandle(c *gin.Context) {
	var req requests.AddUserReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	if query.IsDuplicateLoginNameByUser(req.LoginName) {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.LoginName))
		return
	}
	accountID := handle.GenerateSerialNumber()
	for {
		if !query.IsDuplicateAccountID(accountID) {
			break
		}
	}

	var now = time.Now()
	var accountUser tables.AccountUser
	accountUser.ID = handle.GenerateID()
	accountUser.AccountId = accountID
	accountUser.LoginName = req.LoginName
	accountUser.Password = handle.EncodePassword(accountID, req.Password)
	accountUser.EncodeType = constant.DefaultUserEncryMd5
	accountUser.CreatedAt = now
	accountUser.UpdatedAt = now

	var detail tables.AccountUserDetail
	detail.ID = accountUser.ID

	detail.AccountID = accountID
	detail.Username = req.Username
	detail.CreatedAt = now
	detail.UpdatedAt = now
	detail.IsOfficial = req.IsOfficial

	var beans = make([]interface{}, 0)
	beans = append(beans, &accountUser)
	beans = append(beans, &detail)

	if err := db.GetAccountDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	handle.SuccessResp(c, "", nil)
}

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
	var param = query.UserQueryParam{
		AccountId:    c.Query(handle.QueryAccountId),
		LikeUsername: c.Query(handle.QueryLikeUsername),
		Page:         reqPage.Page,
		PageSize:     reqPage.PageSize,
	}
	details, total, err = query.GetUsers(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
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
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}
