package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"github.com/3115826227/baby-fried-rice/module/account/src/redis"
	"fmt"
)

func RootLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	var root = model.AccountRoot{}
	err = db.DB.Find(&root).Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	token, err := GenerateToken(root.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    root.ID,
		LoginName: root.LoginName,
		Username:  root.Username,
	}
	var loginResult = model.LoginResult{
		UserInfo: userInfo,
		Token:    token,
		Policies: make(map[string][]string),
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var userMeta = &model.UserMeta{
		UserId: root.ID,
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func RootDetail(c *gin.Context) {
	userMeta := GetUserMeta(c)
	root, err := model.GetUserDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}
	var rsp model.RspUserData
	rsp.UserId = root.ID
	rsp.Username = root.Username
	rsp.LoginName = root.LoginName

	SuccessResp(c, "", rsp)
}

func RootLogout(c *gin.Context) {

}
