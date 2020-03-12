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

func AdminLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	var admin = model.AccountAdmin{}
	err = db.DB.Find(&admin).Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	token, err := GenerateToken(admin.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    admin.ID,
		LoginName: admin.LoginName,
		Username:  admin.Username,
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
		UserId: admin.ID,
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func AddAdmin(c *gin.Context) {
	var err error
	var req model.ReqAdminAdd
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	userMeta := GetUserMeta(c)
	_, err = model.GetRootDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}

	var admin = new(model.AccountAdmin)
	admin.ID = GenerateID()
	admin.LoginName = req.LoginName
	admin.Password = EncodePassword(AdminPassword)
	admin.EncodeType = UserEncryMd5

	var beans = make([]interface{}, 0)
	beans = append(beans, &admin)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}
