package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/account/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/log"
	"net/http"
)

func RoleAdd(c *gin.Context) {
	var err error
	var req model.ReqRoleAdd
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}


}
