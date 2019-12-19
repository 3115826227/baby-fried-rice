package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"time"
	"github.com/3115826227/baby-fried-rice/module/account/src/redis"
	"fmt"
	"strings"
)

func UserRegister(c *gin.Context) {
	var err error
	var req model.ReqUserRegister
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if model.IsDuplicateLoginNameByUser(req.LoginName) {
		ErrorResp(c, http.StatusBadRequest, ErrCodeDuplicateName, ErrCodeM[ErrCodeDuplicateName])
		return
	}
	var user model.AccountUser
	user.ID = GenerateID()
	user.LoginName = req.LoginName
	user.Password = EncodePassword(req.Password)
	user.EncodeType = UserEncryMd5
	if err := db.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func UserAdd(c *gin.Context) {

}

func UserLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	var user = model.AccountUser{}
	err = db.DB.Find(&user).Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	token, err := GenerateToken(user.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    user.ID,
		LoginName: user.LoginName,
		Username:  user.Username,
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
		UserId:   user.ID,
		IsSuper:  "0",
		SchoolId: user.SchoolID,
		ReqId:    user.ID,
		Platform: "pc",
	}

	log.Logger.Info(fmt.Sprint("user login header info:",userMeta.ToString()))
	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func UserDetail(c *gin.Context) {
	userMeta := GetUserMeta(c)
	user, err := model.GetUserDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}
	var rsp model.RspUserData
	rsp.UserId = user.ID
	rsp.Username = user.Username
	rsp.SchoolId = user.SchoolID
	rsp.LoginName = user.LoginName

	SuccessResp(c, "", rsp)
}

func UserLogout(c *gin.Context) {

}

func UserUpdate(c *gin.Context) {
	var req model.ReqUserUpdate
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	userMeta := GetUserMeta(c)
	user, err := model.GetUserDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}

	var updateMap = map[string]interface{}{
		"username":    req.Username,
		"update_time": time.Now(),
	}
	if err := db.DB.Model(model.AccountUser{}).Where("id = ?", user.ID).Updates(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func UserDelete(c *gin.Context) {

}
