package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SchoolGet(c *gin.Context) {
	SuccessResp(c, " ", model.GetSchool())
}

func SchoolAdd(c *gin.Context) {
	var req model.ReqSchoolAdd
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var school = &model.School{
		ID:       GenerateID(),
		Name:     req.Name,
		Province: req.Province,
		City:     req.City,
	}
	if err := db.DB.Debug().Create(&school).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}
