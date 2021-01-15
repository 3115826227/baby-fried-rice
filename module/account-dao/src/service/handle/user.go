package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
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
	if IsDuplicateLoginNameByUser(req.LoginName) {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var now = time.Now()
	var user model.AccountUser
	user.ID = GenerateID()
	user.LoginName = req.LoginName
	user.Password = EncodePassword(req.Password)
	user.EncodeType = config.DefaultUserEncryMd5
	user.CreatedAt = now
	user.UpdatedAt = now

	var detail model.AccountUserDetail
	detail.ID = user.ID
	accountID := GenerateSerialNumber()
	for {
		var exist = 0
		if err := db.DB.Debug().Model(&model.AccountUserDetail{}).Where("account_id = ?", accountID).Count(&exist).Error; err != nil {
			continue
		}
		if exist == 0 {
			break
		}
	}
	detail.AccountID = accountID
	detail.Username = req.Username
	detail.Gender = req.Gender
	detail.CreatedAt = now
	detail.UpdatedAt = now

	var userDetail model.UserDetail
	userDetail.UserId = detail.ID
	userDetail.AccountId = detail.AccountID
	userDetail.Username = detail.Username

	var beans = make([]interface{}, 0)
	beans = append(beans, &user)
	beans = append(beans, &detail)
	beans = append(beans, &userDetail)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func UserLogin(c *gin.Context) {
	var err error
	var req model.ReqPasswordLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")

	var user = model.AccountUser{}
	err = db.DB.Debug().Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).First(&user).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	detail := model.GetUserDetail(user.ID)[0]

	go UserLoginLogAdd(user.ID, req.Ip, time.Now())

	var resp = &model.RespUserLogin{
		User:   user,
		Detail: detail,
	}
	SuccessResp(c, "", resp)
}

func UserVerifyDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	user := model.GetUser(id)
	if user.ID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	details := model.GetUserDetail(id)
	if len(details) != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var resp = &model.RespUserDetail{
		User:   user,
		Detail: details[0],
	}
	if details[0].Verify {
		schoolDetail := model.GetUserSchoolDetail(id)
		if len(schoolDetail) != 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
			return
		}
		labels := strings.Split(GetLabel(schoolDetail[0].OrgId), "-")
		var label = model.RespSchoolLabel{
			School:  labels[0],
			Faculty: labels[1],
			Grade:   labels[2],
			Major:   labels[3],
		}
		resp.SchoolLabel = label
		resp.School = schoolDetail[0]
	}
	SuccessResp(c, "", resp)
}

func UserDetail(c *gin.Context) {
	id := c.Query("id")
	var users = make([]model.AccountUserDetail, 0)
	if id != "" {
		users = model.GetUserDetail(id)
	} else {
		users = model.GetUserDetail()
	}
	SuccessResp(c, "", users)
}

func UserSchoolDetail(c *gin.Context) {
	id := c.Query("id")
	var detail = make([]model.AccountUserSchoolDetail, 0)
	if id != "" {
		detail = model.GetUserSchoolDetail(id)
	} else {
		detail = model.GetUserSchoolDetail()
	}
	SuccessResp(c, "", detail)
}

func UserVerify(c *gin.Context) {
	var err error
	var req model.ReqUserVerify
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var student model.AccountSchoolStudent
	var organizes = make([]model.AccountSchoolOrganize, 0)
	if err = db.DB.Debug().Where("school_id = ?", req.SchoolId).Find(&organizes).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var org = make([]string, 0)
	for _, o := range organizes {
		org = append(org, o.ID)
	}
	if err = db.DB.Debug().Where("name = ? and identify = ? and org_id in (?)", req.Name, req.Identify, org).Find(&student).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var updateMap = map[string]interface{}{"school_id": req.SchoolId, "verify": true}
	if err = tx.Model(&model.AccountUserDetail{}).Where("id = ?", req.UserId).Update(updateMap).Error; err != nil {
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
	schoolDetail.ID = req.UserId
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

	SuccessResp(c, "", nil)
}

func UserUpdate(c *gin.Context) {
	var req model.ReqUserUpdate
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	users := model.GetUserDetail(req.UserId)
	if len(users) < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var updateMap = map[string]interface{}{
		"username":    req.Username,
		"update_time": time.Now(),
	}
	if err := db.DB.Model(model.AccountUser{}).Where("id = ?", req.UserId).Updates(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func UserInfoHandle(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var rsp model.RspUserInfo
	var detail model.AccountUserDetail
	details := model.GetUserDetail(id)
	if len(details) != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	detail = details[0]
	rsp.Id = detail.ID
	rsp.Username = detail.Username
	rsp.Verify = detail.Verify
	rsp.Gender = detail.Gender
	rsp.Age = detail.Age

	if detail.Verify {
		school := model.GetUserSchoolDetail(id)[0]
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
