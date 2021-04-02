package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/config"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func UserRegisterHandle(c *gin.Context) {
	var err error
	var req requests.UserRegisterReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	if query.IsDuplicateLoginNameByUser(req.LoginName) {
		log.Logger.Error(fmt.Sprintf("login name %v is duplication", req.LoginName))
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}

	var now = time.Now()
	var user tables.AccountUser
	user.ID = handle.GenerateID()
	user.LoginName = req.LoginName
	user.Password = req.Password
	user.EncodeType = config.DefaultUserEncryMd5
	user.CreatedAt = now
	user.UpdatedAt = now

	var detail tables.AccountUserDetail
	detail.ID = user.ID
	accountID := handle.GenerateSerialNumber()
	for {
		if !query.IsDuplicateAccountID(accountID) {
			break
		}
	}

	detail.AccountID = accountID
	detail.Username = req.Username
	detail.Gender = req.Gender
	detail.Phone = req.Phone
	detail.CreatedAt = now
	detail.UpdatedAt = now

	var userDetail tables.UserDetail
	userDetail.UserId = detail.ID
	userDetail.AccountId = detail.AccountID
	userDetail.Username = detail.Username

	var beans = make([]interface{}, 0)
	beans = append(beans, &user)
	beans = append(beans, &detail)
	beans = append(beans, &userDetail)

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}

	handle.SuccessResp(c, "", nil)
}

func UserLoginHandle(c *gin.Context) {
	var err error
	var req requests.PasswordLoginReq
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")

	user, err := query.GetUserByLogin(req.LoginName, req.Password)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detail = new(tables.AccountUserDetail)
	detail.ID = user.ID
	if err = db.GetDB().GetObject(nil, detail); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp = model.RespUserLogin{
		User:   user,
		Detail: *detail,
	}
	handle.SuccessResp(c, "", resp)
}

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
