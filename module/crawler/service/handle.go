package service

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/3115826227/baby-fried-rice/module/crawler/model/db"
	"github.com/3115826227/baby-fried-rice/module/crawler/redis"
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func findTrainMeta(date string, isDetail bool) (rsp model.RspTrainMetaTrigger) {
	rsp.Date = date
	key := fmt.Sprintf("%v:%v", TrainMetaMember, date)
	if redis.Exist(key) {
		rsp.IsTrigger = true
		var trainMeta = make([]model.TrainMeta, 0)
		if err := db.DB.Where("date = ?", date).Find(&trainMeta).Error; err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		rsp.TriggerTrainNumber = len(trainMeta)
		if isDetail {
			var triggerTrains = make([]model.RspTriggerTrain, 0)
			for _, train := range trainMeta {
				triggerTrains = append(triggerTrains, model.RspTriggerTrain{
					Train:         train.Train,
					StartStation:  train.StartStation,
					StartTime:     train.StartTime,
					ArriveStation: train.ArriveStation,
					ArriveTime:    train.ArriveTime,
					OverDay:       train.OverDay,
				})
			}
			rsp.TriggerTrain = triggerTrains
		}
	}
	return
}

func FindTrainMeta(c *gin.Context) {
	date := c.Query("date")
	var resp = make([]model.RspTrainMetaTrigger, 0)
	if date == "" {
		today := time.Now().Format(config.DayLayout)
		thirtyDay := time.Now().AddDate(0, 0, 30).Format(config.DayLayout)
		tempDate, _ := time.Parse(config.DayLayout, today)
		for {
			var tempDateStr = tempDate.Format(config.DayLayout)
			resp = append(resp, findTrainMeta(tempDateStr, false))
			if tempDateStr == thirtyDay {
				break
			}
			tempDate = tempDate.AddDate(0, 0, 1)
		}
	} else if IsValidDate(date) {
		resp = append(resp, findTrainMeta(date, true))
	} else {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	SuccessResp(c, "", resp)
}

func FindTrainSeat(c *gin.Context) {

}
