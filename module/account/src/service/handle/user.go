package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/redis"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
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

	var detail model.AccountUserDetail
	detail.ID = user.ID

	var beans = make([]interface{}, 0)
	beans = append(beans, &user)
	beans = append(beans, &detail)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	用户实名制认证
*/
func UserVerify(c *gin.Context) {
	var err error
	var req model.ReqUserVerify
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var certify model.SchoolUserCertification
	certify, err = model.FindCertification(req.Identify, req.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var department model.SchoolDepartment
	department, err = model.FindDepartmentById(certify.SchoolDepartmentId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	userMeta := GetUserMeta(c)
	tx := db.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var updateMap = map[string]interface{}{"school_id": department.SchoolId, "verify": 1}
	if err = tx.Model(model.AccountUserDetail{}).Where("id = ?", userMeta.UserId).Update(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var schoolDetail model.AccountUserSchoolDetail
	schoolDetail.ID = userMeta.UserId
	schoolDetail.Name = certify.Name
	schoolDetail.Identify = certify.Identify
	schoolDetail.SchoolDepartmentId = certify.SchoolDepartmentId
	if err = tx.Create(&schoolDetail).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", model.RspOkResponse{})
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
	err = db.DB.Debug().Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Find(&user).Error
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

	var detail = model.AccountUserDetail{}
	err = db.DB.Where("id = ?", user.ID).First(&detail).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:   user.ID,
		Username: detail.Username,
		SchoolId: detail.SchoolId,
	}
	var loginResult = model.LoginResult{
		UserInfo:   userInfo,
		Token:      token,
		Permission: make([]int, 0),
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var userMeta = &model.UserMeta{
		UserId:   user.ID,
		IsSuper:  "0",
		SchoolId: detail.SchoolId,
		ReqId:    user.ID,
		Platform: "pc",
	}

	log.Logger.Info(fmt.Sprint("user login header info:", userMeta.ToString()))
	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
	redis.AddAccountToken(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))

	c.JSON(http.StatusOK, result)
}

/*
	用户退出登录
*/
func UserLogout(c *gin.Context) {
	userMeta := GetUserMeta(c)
	token, err := redis.GetAccountToken(userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	redis.DeleteAccountToken(token)
	redis.DeleteAccountToken(userMeta.UserId)
	SuccessResp(c, "", model.RspOkResponse{})
}

/*
	用户登录token刷新
*/
func UserRefresh(c *gin.Context) {
	userMeta := GetUserMeta(c)
	token, err := GenerateToken(userMeta.UserId, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	oldToken, err := redis.GetAccountToken(userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	//删除旧token，更新新的token
	redis.DeleteAccountToken(oldToken)
	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
	redis.AddAccountToken(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))
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
	rsp.SchoolId = user.SchoolId

	SuccessResp(c, "", rsp)
}

func UserPasswordUpdate(c *gin.Context) {

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
