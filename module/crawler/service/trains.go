package service

import (
	"fmt"
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/json-iterator/go"
	"github.com/3115826227/baby-fried-rice/module/crawler/redis"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/3115826227/baby-fried-rice/module/crawler/model/db"
	"go.uber.org/zap"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/crawler/utils"
	"strings"
	"strconv"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
)

const (
	trainTaskKey = "train:task"
)

type TrainTask struct {
	Retry int
	Train model.TrainMeta
}

func (task *TrainTask) ToString() string {
	data, _ := jsoniter.Marshal(task)
	return string(data)
}

type TrainRelationResult struct {
	train         model.TrainMeta
	relationBeans [][]interface{}
}

type CityResult struct {
	Station     string
	StationCode string
	City        string
}

type TrainConsumerClose struct {
	date string
}

//列车爬取任务
var trainTaskChan = make(chan TrainTask, 5000)
//列车爬取基本信息
var trainMetaChan = make(chan model.TrainMeta, 500)
//列车过站信息
var trainRelationChan = make(chan []model.TrainStationRelation, 500)

func TrainMetaTrigger(c *gin.Context) {
	date := c.Query("date")
	if date == "" || !IsValidDate(date) {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	go TrainMetaExecutor(date)

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func MeituanTrainConsumer() {
	tick := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-tick:
			num = 0
		default:
			if num > 3 {
				continue
			}
			task := <-trainTaskChan
			go meituanTrainTrigger(task)
			num += 1
		}
	}
}

func meituanTrainTrigger(task TrainTask) {
	var err error
	defer func() {
		if err != nil {
			if task.Retry < 2 {
				task.Retry += 1
				redis.LPush(trainTaskKey, task.ToString())
			}
		}
	}()

	train := task.Train
	URL := fmt.Sprintf("https://i.meituan.com/uts/train/train/timetable?fromPC=1&train_source=meituanpc@wap&train_code=%v&start_date=%v&from_station_telecode=%v&to_station_telecode=%v",
		train.Train, train.Date, redis.Get(train.StartStation), redis.Get(train.ArriveStation))

	data, err := utils.Request(URL)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var response model.RspMeituanTrainMeta
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.Data.Stations) == 0 {
		if task.Retry < 2 {
			task.Retry += 1
			redis.LPush(trainTaskKey, task.ToString())
		}
		return
	}

	list := make([]model.TrainStationRelation, 0)
	var num = 1
	for _, item := range response.Data.Stations {
		stationCode := redis.Get(item.StationName)
		var stopTime = 0
		if num == 1 {
			item.ArriveTime = item.StartTime
		} else if num == len(response.Data.Stations) {
			item.StartTime = item.ArriveTime
		} else {
			stopTime, err = strconv.Atoi(item.StopTime)
			if err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
		}

		list = append(list, model.TrainStationRelation{
			Train:         train.Train,
			Date:          train.Date,
			StationNumber: num,
			Station:       item.StationName,
			StationCode:   stationCode,
			StartTime:     item.StartTime,
			ArriveTime:    item.ArriveTime,
			StopMinute:    stopTime,
		})
		num += 1
	}

	var trainMeta = model.TrainMeta{
		Train:                train.Train,
		TrainCode:            train.TrainCode,
		Date:                 train.Date,
		StartStation:         train.StartStation,
		StartStationCode:     redis.Get(train.StartStation),
		StartTime:            list[0].StartTime,
		ArriveStation:        train.ArriveStation,
		ArriveStationCode:    redis.Get(train.ArriveStation),
		ArriveTime:           list[len(list)-1].ArriveTime,
		RunningStationNumber: len(list),
	}

	var compute = &TrainComputeInfo{Train: trainMeta, Stations: list}
	compute.ComputeRunningAndOverDay()
	trainRelationChan <- compute.Stations
	trainMetaChan <- compute.Train
}

func ZhixingTrainConsumer() {
	tick := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-tick:
			num = 0
		default:
			if num > 10 {
				continue
			}
			task := <-trainTaskChan
			go zhixingTrainTrigger(task)
			num += 1
		}
	}
}

func zhixingTrainTrigger(task TrainTask) {
	var err error
	defer func() {
		if err != nil {
			if task.Retry < 2 {
				task.Retry += 1
				redis.LPush(trainTaskKey, task.ToString())
			}
		}
	}()

	train := task.Train
	trafficBaseZhiXing := "http://m.ctrip.com/restapi/soa2/10103/json/GetTrainStopListV3"

	payload := strings.NewReader(fmt.Sprintf("{\"DepartStation\":\"%v\",\"ArriveStation\":\"%v\",\"DepartureDate\":\"%v\",\"TrainName\":\"%v\"}",
		train.StartStation, train.ArriveStation, train.Date, train.Train))

	data, err := utils.PostRequest("POST", trafficBaseZhiXing, payload)
	if err != nil {
		return
	}

	var response model.RspZhiXingTrainMeta
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.TrainStopList) == 0 {
		if task.Retry < 2 {
			task.Retry += 1
			redis.LPush(trainTaskKey, task.ToString())
		}
		return
	}

	list := make([]model.TrainStationRelation, 0)
	var num = 1
	for _, item := range response.TrainStopList {
		stationCode := redis.Get(item.StationName)
		if num == 1 {
			item.ArrivalTime = item.DepartureTime
		}
		if num == len(response.TrainStopList) {
			item.DepartureTime = item.ArrivalTime
		}
		list = append(list, model.TrainStationRelation{
			Train:         train.Train,
			Date:          train.Date,
			StationNumber: num,
			Station:       item.StationName,
			StationCode:   stationCode,
			StartTime:     item.DepartureTime,
			ArriveTime:    item.ArrivalTime,
			StopMinute:    item.StopTimes,
		})
		num += 1
	}

	var trainMeta = model.TrainMeta{
		Train:                train.Train,
		TrainCode:            train.TrainCode,
		Date:                 train.Date,
		StartStation:         train.StartStation,
		StartStationCode:     redis.Get(train.StartStation),
		StartTime:            list[0].StartTime,
		ArriveStation:        train.ArriveStation,
		ArriveStationCode:    redis.Get(train.ArriveStation),
		ArriveTime:           list[len(list)-1].ArriveTime,
		RunningStationNumber: len(list),
	}

	var compute = &TrainComputeInfo{Train: trainMeta, Stations: list}
	compute.ComputeRunningAndOverDay()
	trainRelationChan <- compute.Stations
	trainMetaChan <- compute.Train
}

func QunarTrainConsumer() {
	tick := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-tick:
			num = 0
		default:
			if num > 0 {
				continue
			}
			task := <-trainTaskChan
			go QunarTrainTrigger(task)
			num += 1
		}
	}
}

func QunarTrainTrigger(task TrainTask) {
	var err error
	defer func() {
		if err != nil {
			if task.Retry < 2 {
				task.Retry += 1
				redis.LPush(trainTaskKey, task.ToString())
			}
		}
	}()

	train := task.Train
	URL := fmt.Sprintf("https://train.qunar.com/dict/open/seatDetail.do?dptStation=%v&arrStation=%v&date=%v&trainNo=%v&needTimeDetail=true",
		train.StartStation, train.ArriveStation, train.Date, train.Train)
	data, err := utils.Request(URL)
	if err != nil {
		return
	}

	var response model.RspTrainMetaWayStation
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.Data.StationItemList) == 0 {
		if task.Retry < 2 {
			task.Retry += 1
			redis.LPush(trainTaskKey, task.ToString())
		}
		return
	}

	list := make([]model.TrainStationRelation, 0)
	var num = 1
	for _, item := range response.Data.StationItemList {
		stationCode := redis.Get(item.StationName)
		if num == 1 {
			item.ArriveTime = item.StartTime
		}
		if num == len(response.Data.StationItemList) {
			item.StartTime = item.ArriveTime
		}
		list = append(list, model.TrainStationRelation{
			Train:         train.Train,
			Date:          train.Date,
			StationNumber: num,
			Station:       item.StationName,
			StationCode:   stationCode,
			StartTime:     item.StartTime,
			ArriveTime:    item.ArriveTime,
			StopMinute:    item.OverTime,
		})
		num += 1
	}

	var trainMeta = model.TrainMeta{
		Train:                train.Train,
		TrainCode:            train.TrainCode,
		Date:                 train.Date,
		StartStation:         train.StartStation,
		StartStationCode:     redis.Get(train.StartStation),
		StartTime:            list[0].StartTime,
		ArriveStation:        train.ArriveStation,
		ArriveStationCode:    redis.Get(train.ArriveStation),
		ArriveTime:           list[len(list)-1].ArriveTime,
		RunningStationNumber: len(list),
	}

	var compute = &TrainComputeInfo{Train: trainMeta, Stations: list}
	compute.ComputeRunningAndOverDay()
	trainRelationChan <- compute.Stations
	trainMetaChan <- compute.Train
}

func TrainMetaExecutor(date string) {
	stations := GetStation()
	trainMetaStationTriggerKey := fmt.Sprintf("%v:%v", config.TrainMetaStationTriggerPrefix, date)
	trainMetaCodeTriggerKey := fmt.Sprintf("%v:%v", config.TrainMetaCodeTriggerPrefix, date)
	for _, station := range stations {
		if redis.HashExist(trainMetaStationTriggerKey, station.Name) {
			continue
		}

		var URL = fmt.Sprintf("https://kyfw.12306.cn/otn/czxx/query?train_start_date=%v&train_station_name=%v&train_station_code=%v&randCode=",
			date, station.Name, station.Code)
		data, err := utils.Request(URL)
		if err != nil {
			continue
		}

		var response model.RspStationTrainMeta
		err = jsoniter.Unmarshal(data, &response)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}

		for _, item := range response.Data.Data {
			if redis.HashExist(trainMetaCodeTriggerKey, item.TrainNo) {
				continue
			}
			var train = model.TrainMeta{
				Date:          date,
				Train:         item.StationTrainCode,
				TrainCode:     item.TrainNo,
				StartStation:  item.StartStationName,
				ArriveStation: item.EndStationName,
			}

			task := TrainTask{Train: train, Retry: 0}
			redis.LPush(trainTaskKey, task.ToString())
			redis.HashAdd(trainMetaCodeTriggerKey, item.TrainNo, time.Now().String())
		}

		redis.HashAdd(trainMetaStationTriggerKey, station.Name, time.Now().String())
	}
}

func TrainTaskConsumer() {
	for {
		result := redis.BRPop(trainTaskKey)
		if len(result) == 0 {
			continue
		}
		var info = TrainTask{}
		err := jsoniter.Unmarshal([]byte(result[1]), &info)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		trainTaskChan <- info
	}
}

func TrainMetaInsertConsumer() {
	tick := time.Tick(5 * time.Second)
	var list = make([]model.TrainMeta, 0)
	for {
		select {
		case <-tick:
			var beans = make([]model.TrainMeta, 0)
			for _, train := range list {
				beans = append(beans, train)
			}
			batchInsertTrainMeta(beans)
			list = make([]model.TrainMeta, 0)
		default:
			train := <-trainMetaChan
			list = append(list, train)
		}
	}
}

func TrainRelInsertConsumer() {
	tick := time.Tick(5 * time.Second)
	var list = make([]model.TrainStationRelation, 0)
	for {
		select {
		case <-tick:
			var beans = make([]model.TrainStationRelation, 0)
			for _, relation := range list {
				beans = append(beans, relation)
			}
			batchInsertTrainStationRelation(beans)
			list = make([]model.TrainStationRelation, 0)
		default:
			relations := <-trainRelationChan
			list = append(list, relations...)
		}
	}
}

func batchInsertTrainStationRelation(list []model.TrainStationRelation) (count int, err error) {
	var beans = make([][]interface{}, 0)

	for _, item := range list {
		beans = append(beans, []interface{}{item.Train, item.Date, item.Station, item.StationCode, item.StationNumber,
			item.ArriveTime, item.StartTime, item.StopMinute, item.OverDay})
	}

	var uniqueFields = []string{"train", "date", "station_number"}
	var allFields = []string{"train", "date", "station", "station_code", "station_number", "arrive_time", "start_time", "stop_minute", "over_day"}
	count, err = db.Load((&model.TrainStationRelation{}).TableName(), allFields, uniqueFields, nil, beans)
	if err != nil {
		log.Logger.Error("load to db error")
	} else {
		log.Logger.Info("load to db", zap.Int("records", count))
	}
	return
}

func batchInsertTrainMeta(list []model.TrainMeta) (count int, err error) {
	var beans = make([][]interface{}, 0)

	for _, item := range list {
		beans = append(beans, []interface{}{item.Train, item.TrainCode, item.Date, item.StartStation, item.StartStationCode,
			item.StartTime, item.ArriveStation, item.ArriveStationCode, item.ArriveTime, item.RunningStationNumber,
			item.RunningMinute, item.OverDay})
	}

	var uniqueFields = []string{"train", "date"}
	var allFields = []string{"train", "train_code", "date", "start_station", "start_station_code", "start_time",
		"arrive_station", "arrive_station_code", "arrive_time", "running_station_number", "running_minute", "over_day"}

	count, err = db.Load((&model.TrainMeta{}).TableName(), allFields, uniqueFields, nil, beans)
	if err != nil {
		log.Logger.Error("load to db error")
	} else {
		log.Logger.Info("load to db", zap.Int("records", count))
	}
	return
}
