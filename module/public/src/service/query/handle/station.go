package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func CityGet(c *gin.Context) {
	var stations = make([]model.Station, 0)
	if err := db.DB.Find(&stations).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	sort.Sort(model.Stations(stations))
	var rsp = make([]model.RspStationCity, 0)
	var mp = make(map[string]struct{})
	for _, station := range stations {
		if station.City == "" {
			continue
		}
		if _, exist := mp[station.City]; exist {
			continue
		}
		rsp = append(rsp, model.RspStationCity{Name: station.City})
		mp[station.City] = struct{}{}
	}

	SuccessResp(c, "", rsp)
}
