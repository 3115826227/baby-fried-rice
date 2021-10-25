package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/cache"
	"baby-fried-rice/internal/pkg/module/manage/config"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 超级用户账号登录
func AdminLoginHandle(c *gin.Context) {
	var req requests.PasswordLoginReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	// 根据登录账号查询用户
	req.LoginName = strings.TrimSpace(req.LoginName)
	admin, err := query.GetAdminByLogin(req.LoginName)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 根据用户id和原始密码进行加密
	req.Password = handle.EncodePassword(admin.ID, strings.TrimSpace(req.Password))
	if req.Password != admin.Password {
		err = errors.New("password is invalid")
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	req.Ip = c.GetHeader("IP")

	var token string
	token, err = handle.GenerateToken(admin.ID, time.Now(), config.GetConfig().TokenSecret)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var response = rsp.LoginResult{
		UserInfo: rsp.UserDataResp{
			UserId:    admin.ID,
			LoginName: admin.LoginName,
			Username:  admin.Username,
		},
		Token: token,
	}

	go func() {
		// 写入缓存
		var userMeta = &handle.UserMeta{
			AccountId: admin.ID,
			Username:  admin.Username,
			Platform:  "pc",
		}
		cache.GetCache().Add(fmt.Sprintf("%v:%v", constant.TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.AccountId, fmt.Sprintf("%v:%v", constant.TokenPrefix, token))
		// 写入日志
		var loginLog = tables.AccountAdminLoginLog{
			AccountId: admin.ID,
			IP:        req.Ip,
			LoginTime: time.Now(),
		}
		if err = db.GetAccountDB().CreateObject(&loginLog); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}()

	handle.SuccessResp(c, "", response)
}

func AdminLogoutHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	token, err := cache.GetCache().Get(userMeta.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		cache.GetCache().Del(token)
		cache.GetCache().Del(userMeta.AccountId)
	}()

	handle.SuccessResp(c, "", nil)
}
