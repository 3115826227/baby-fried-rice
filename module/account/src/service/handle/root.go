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

func RootAdd() error {
	var count = 0
	err := db.DB.Model(&model.AccountRoot{}).Count(&count).Error
	if err != nil {
		log.Logger.Warn(err.Error())
		return err
	}
	if count > 0 {
		log.Logger.Info("超级管理账号已存在")
		return nil
	}
	var root = new(model.AccountRoot)
	root.ID = GenerateID()
	root.LoginName = "root"
	root.Password = EncodePassword("root")
	root.EncodeType = UserEncryMd5

	var beans = make([]interface{}, 0)
	beans = append(beans, &root)
	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
		return err
	}
	log.Logger.Info("初始化Root账号")
	return nil
}

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
	err = db.DB.Debug().Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Find(&root).Error
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
		UserInfo:   userInfo,
		Token:      token,
		Permission: make([]int, 0),
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var userMeta = &model.UserMeta{
		UserId:   root.ID,
		IsSuper:  "1",
		SchoolId: "",
		ReqId:    root.ID,
		Platform: "pc",
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func RootDetail(c *gin.Context) {
	userMeta := GetUserMeta(c)
	root, err := model.GetRootDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}
	var rsp model.RspUserData
	rsp.UserId = root.ID

	SuccessResp(c, "", rsp)
}

func RootLogout(c *gin.Context) {

}
