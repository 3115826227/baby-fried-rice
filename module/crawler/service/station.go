package service

import (
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/3115826227/baby-fried-rice/module/crawler/model/db"
	"github.com/3115826227/baby-fried-rice/module/crawler/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func StationTrigger(c *gin.Context) {
	Station()
	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func Station() {
	var URL = "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js"
	resp, err := http.Get(URL)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	str := strings.Split(string(data), "=")[1][1:]
	str = str[:len(str)-2]

	StrList := strings.Split(str, "@")[1:]

	var list = make([]model.Station, 0)
	for _, info := range StrList {
		infoList := strings.Split(info, "|")
		name := infoList[1]
		code := infoList[2]
		now := time.Now()
		var station = new(model.Station)
		station.Name = strings.Replace(name, " ", "", -1)
		station.Code = code
		station.CreatedAt, station.UpdatedAt = now, now
		list = append(list, *station)
		redis.HashAdd(config.StationCodeNameKey, station.Code, station.Name)
		redis.HashAdd(config.StationNameCodeKey, station.Name, station.Code)
	}

	beans := make([][]interface{}, 0)
	uniqueFields := []string{"name", "code"}
	allFields := []string{"name", "code", "create_time", "update_time"}
	for _, station := range list {
		beans = append(beans, []interface{}{station.Name, station.Code, station.CreatedAt, station.UpdatedAt})
	}
	count, err := db.Load((&model.Station{}).TableName(), allFields, uniqueFields, nil, beans)
	if err != nil {
		log.Logger.Error("load to db error")
		return
	} else {
		log.Logger.Info("load to db", zap.Int("records", count))
	}
}
