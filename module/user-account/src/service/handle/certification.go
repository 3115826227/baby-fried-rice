package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"fmt"
)

func SchoolCertificationDelete(c *gin.Context) {
	var err error
	var req model.ReqSchoolCertificationDelete
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	fmt.Println(req)

}
