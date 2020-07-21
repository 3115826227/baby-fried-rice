package handle

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SchoolGet(c *gin.Context) {
	var schools = make([]model.School, 0)
	if err := db.DB.Debug().Model(&model.School{}).Find(&schools).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	SuccessResp(c, "", schools)
}
