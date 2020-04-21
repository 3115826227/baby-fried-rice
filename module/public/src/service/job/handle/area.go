package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func areaGetByCity(cityCode string) (result model.RspArea, err error) {
	var city model.Area
	if err = db.DB.Where("code = ?", cityCode).Find(&city).Error; err != nil {
		log.Logger.Warn(err.Error())
		return result, err
	}
	var province model.Area
	if err = db.DB.Where("parent_code = ?", city.ParentCode).Find(&province).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var locals = make([]model.Area, 0)
	if err = db.DB.Where("parent_code = ?", city.Code).Find(&locals).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var rspLocals = make([]model.RspLocal, 0)
	for _, local := range locals {
		rspLocals = append(rspLocals, model.RspLocal{Local: local.Name, Code: local.Code})
	}
	result.Province = province.Name
	result.Code = province.Code
	result.Cities = []model.RspCity{{
		City:   city.Name,
		Code:   city.Code,
		Locals: rspLocals,
	}}
	return
}

func AreaGet(c *gin.Context) {
	provinceCode := c.Query("province")
	cityCode := c.Query("city")

	var result = make([]model.RspArea, 0)
	if provinceCode == "" && cityCode == "" {
		//
	} else if provinceCode == "" {
		area, err := areaGetByCity(cityCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
		result = append(result, area)
	} else if cityCode == "" {

	} else {

	}
}
