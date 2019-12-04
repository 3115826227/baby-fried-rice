package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"sort"
)

func stationGet() (cities []model.StationCity, err error) {
	var stations = make([]model.Station, 0)
	if err = db.DB.Group("city").Find(stations).Error; err != nil {
		log.Logger.Warn(err.Error())
		return cities, err
	}
	for _, stations := range stations {
		cities = append(cities, model.StationCity{City: stations.City})
	}
	sort.Sort(model.StationCities(cities))
	return
}

func StationGet(c *gin.Context) {

}
