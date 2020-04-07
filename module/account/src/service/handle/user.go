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
	var now = time.Now()
	var user model.AccountUser
	user.ID = GenerateID()
	user.LoginName = req.LoginName
	user.Password = EncodePassword(req.Password)
	user.EncodeType = UserEncryMd5
	user.CreatedAt = now
	user.UpdatedAt = now

	var detail model.AccountUserDetail
	detail.ID = user.ID
	detail.Username = req.Username
	detail.Gender = req.Gender
	detail.CreatedAt = now
	detail.UpdatedAt = now

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

	var student model.AccountSchoolStudent
	var organizes = make([]model.AccountSchoolOrganize, 0)
	if err := db.DB.Debug().Where("school_id = ?", req.SchoolId).Find(&organizes).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var org = make([]string, 0)
	for _, o := range organizes {
		org = append(org, o.ID)
	}
	if err := db.DB.Debug().Where("name = ? and identify = ? and org_id in (?)", req.Name, req.Identify, org).Find(&student).Error; err != nil {
		log.Logger.Warn(err.Error())
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

	var updateMap = map[string]interface{}{"school_id": req.SchoolId, "verify": true}
	if err = tx.Model(&model.AccountUserDetail{}).Where("id = ?", userMeta.UserId).Update(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var updateStudentMap = map[string]interface{}{"status": true}
	if err = tx.Model(&model.AccountSchoolStudent{}).Where("id = ?", student.ID).Update(updateStudentMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var schoolDetail model.AccountUserSchoolDetail
	schoolDetail.ID = userMeta.UserId
	schoolDetail.Name = student.Name
	schoolDetail.Identify = student.Identify
	schoolDetail.Number = student.Number
	schoolDetail.OrgId = student.OrgId
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
		Username: detail.Username,
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
	var rsp model.RspUserDetail
	var login = model.AccountUser{}
	if err := db.DB.Debug().Where("id = ?", userMeta.UserId).Find(&login).Error; err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}
	rsp.LoginName = login.LoginName
	rsp.UserId = user.ID
	rsp.Username = user.Username
	rsp.SchoolId = user.SchoolId
	rsp.Gender = user.Gender
	rsp.Age = user.Age
	rsp.Verify = user.Verify
	if user.Verify {
		detail, err := model.GetUserSchoolDetail(user.ID)
		if err != nil {
			ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
			return
		}
		labels := strings.Split(GetLabel(detail.OrgId), "-")
		rsp.Faculty = labels[1]
		rsp.Grade = labels[2]
		rsp.Major = labels[3]
		var school model.School
		if err := db.DB.Debug().Model(&model.School{}).Where("id = ?", user.SchoolId).Find(&school).Error; err != nil {
			log.Logger.Warn(err.Error())
			ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
			return
		}
		rsp.School = school.Name
		rsp.Name = detail.Name
		rsp.Identify = detail.Identify
		rsp.Number = detail.Number
	}

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

func UserInfoGet(c *gin.Context) {
	id := c.Query("id")

	var rsp model.RspUserInfo
	var detail model.AccountUserDetail
	if err := db.DB.Debug().Model(&model.AccountUserDetail{}).Where("id = ?", id).Find(&detail).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	rsp.Id = detail.ID
	rsp.Username = detail.Username
	rsp.Verify = detail.Verify
	rsp.Gender = detail.Gender
	rsp.Age = detail.Age

	if detail.Verify {
		var school model.AccountUserSchoolDetail
		if err := db.DB.Debug().Model(&model.AccountUserSchoolDetail{}).Where("id = ?", id).Find(&school).Error; err != nil {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
		labels := strings.Split(GetLabel(school.OrgId), "-")
		rsp.School = labels[0]
		rsp.Faculty = labels[1]
		rsp.Grade = labels[2]
		rsp.Major = labels[3]
		rsp.Name = school.Name
		rsp.Number = school.Number
	}

	SuccessResp(c, "", rsp)
}
