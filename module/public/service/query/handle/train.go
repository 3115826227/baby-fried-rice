package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"github.com/3115826227/baby-fried-rice/module/public/log"
)

func TrainMetaGet(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	var train model.TrainMeta
	if err := db.DB.Where("train = ?", code).Find(&train).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", train)
}
