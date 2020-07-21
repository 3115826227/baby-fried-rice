package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func RootLogin(c *gin.Context) {
	var err error
	var req model.ReqPasswordLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")

	var root = model.AccountRoot{}
	err = db.DB.Debug().Where("login_name = ? and password = ?", req.LoginName, EncodePassword(req.Password)).Find(&root).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	go RootLoginLogAdd(root.ID, req.Ip, time.Now())

	SuccessResp(c, "", root)
}

func RootDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	SuccessResp(c, "", model.GetRoot(id))
}
