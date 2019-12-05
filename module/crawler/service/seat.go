package service

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/3115826227/baby-fried-rice/module/crawler/model/db"
	"github.com/3115826227/baby-fried-rice/module/crawler/redis"
	"github.com/3115826227/baby-fried-rice/module/crawler/utils"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
	"gopkg.in/gin-gonic/gin.v1/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	seatConsumerKey = "seat"
)

type SeatRelationInfo struct {
	Train             string `json:"train"`
	Date              string `json:"date"`
	From              string `json:"from"`
	FromStationNumber int    `json:"from_station_number"`
	To                string `json:"to"`
	ToStationNumber   int    `json:"to_station_number"`
}

func (info *SeatRelationInfo) ToString() string {
	data, _ := json.Marshal(info)
	return string(data)
}

type SeatPriceInfo struct {
	Date  string                        `json:"date"`
	From  string                        `json:"from"`
	To    string                        `json:"to"`
	Beans []model.TrainStationSeatPrice `json:"beans"`
}

func (info *SeatPriceInfo) ToString() string {
	data, _ := json.Marshal(info)
	return string(data)
}

var seatRelationChan = make(chan SeatRelationInfo, 5000)
var seatPriceChan = make(chan SeatPriceInfo, 5000)

func IsTrainSeatTrigger(date string) bool {
	status := redis.HashGet(config.TrainSeatTriggerStatus, date)
	if status == config.RunningStatus || status == config.SuccessStatus {
		return false
	}
	return true
}

func UpdateTrainSeatTriggerStatus(date, status string) {
	redis.HashAdd(config.TrainSeatTriggerStatus, date, status)
}

func TrainSeatTrigger(c *gin.Context) {
	date := c.Query("date")
	if date == "" || !IsValidDate(date) {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	//if !IsTrainSeatTrigger(date) {
	//	ErrorResp(c, http.StatusInternalServerError, 203, "信息正在抓取或者已经抓取成功")
	//}

	//UpdateTrainSeatTriggerStatus(date, config.RunningStatus)
	go TrainSeatPrice(date)

	c.JSON(http.StatusOK, model.RspOkResponse{Message: "信息已经开始抓取"})

}

func JindongConsumer() {
	ticker := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-ticker:
			num = 0
		default:
			if num > 5 {
				continue
			}
			info := <-seatRelationChan
			go jindongTraffic(info)
			num += 1
		}
	}
}

func jindongTraffic(info SeatRelationInfo) {
	var err error
	defer func() {
		if err != nil {
			redis.LPush(seatConsumerKey, info.ToString())
		}
	}()

	URL := "https://train.jd.com/query/getTrainTickets.html"
	payload := strings.NewReader(fmt.Sprintf("stationQuery.fromStation=%v&stationQuery.toStation=%v&stationQuery.date=%v", info.From, info.To, info.Date))
	req, err := http.NewRequest("POST", URL, payload)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	req.Header.Add("referer", fmt.Sprintf("https://train.jd.com/query/queryTrains.html?stationQuery.fromStation=%v&stationQuery.toStation=%v&stationQuery.date=%v&stationQuery.requestType=0", info.From, info.To, info.Date))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var response model.RspJindongTrainSeatPrice
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.Data.Value) == 0 {
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}
	var list = make([]model.TrainStationSeatPrice, 0)
	for _, item := range response.Data.Value {
		for _, seats := range item.Seats {
			price, err := strconv.ParseFloat(seats.Price, 64)
			if err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			list = append(list, model.TrainStationSeatPrice{
				Train:             item.TrainCode,
				Date:              info.Date,
				StartStation:      info.From,
				StartStationCode:  redis.Get(info.From),
				ArriveStation:     info.To,
				ArriveStationCode: redis.Get(info.To),
				SeatCategoryName:  seats.SeatName,
				Price:             int(price * 100),
			})
		}
	}

	seatPriceChan <- SeatPriceInfo{Date: info.Date, From: info.From, To: info.To, Beans: list}
}

func MeituanConsumer() {
	ticker := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-ticker:
			num = 0
		default:
			if num > 5 {
				continue
			}
			info := <-seatRelationChan
			go meituanTraffic(info)
			num += 1
		}
	}
}

func meituanTraffic(info SeatRelationInfo) {
	var err error
	defer func() {
		if err != nil {
			redis.LPush(seatConsumerKey, info.ToString())
		}
	}()

	URL := fmt.Sprintf("https://i.meituan.com/uts/train/train/querytripnew?fromPC=1&train_source=meituanpc@wap&from_station_telecode=%v&to_station_telecode=%v&start_date=%v&isStudentBuying=false",
		redis.Get(info.From), redis.Get(info.To), info.Date)

	data, err := utils.Request(URL)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var response model.RspMeituanTrainSeatPrice
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.Data.Trains) == 0 {
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}

	var list = make([]model.TrainStationSeatPrice, 0)
	for _, item := range response.Data.Trains {
		for _, seats := range item.Seats {
			list = append(list, model.TrainStationSeatPrice{
				Train:             item.TrainCode,
				Date:              info.Date,
				StartStation:      info.From,
				StartStationCode:  redis.Get(info.From),
				ArriveStation:     info.To,
				ArriveStationCode: redis.Get(info.To),
				SeatCategoryName:  seats.SeatTypeName,
				Price:             int(seats.SeatPrice * 100),
			})
		}
	}

	cityUpdate(info.From, response.Data.FromCityName)
	cityUpdate(info.To, response.Data.ToCityName)
	seatPriceChan <- SeatPriceInfo{Date: info.Date, From: info.From, To: info.To, Beans: list}
}

func TongChengYiLongConsumer() {
	ticker := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-ticker:
			num = 0
		default:
			if num > 0 {
				continue
			}
			info := <-seatRelationChan
			go tongChengYiLongTraffic(info)
			num += 1
		}
	}
}

func tongChengYiLongTraffic(info SeatRelationInfo) {
	var err error
	defer func() {
		if err != nil {
			redis.LPush(seatConsumerKey, info.ToString())
		}
	}()

	trafficBaseTongChenYiLong := `https://www.ly.com/uniontrain/trainapi/TrainPCCommon/SearchTrainRemainderTickets`

	v := url.Values{}
	v.Add("para", fmt.Sprintf(`{"To":"%v","From":"%v","TrainDate":"%v","constId":"4yDOvPn_ETzqy_vnk4eVJ9zOSL3OLuc-Exmux_tymao","platId":1,"headver":"1.0.0"}`,
		info.To, info.From, info.Date))

	URL := fmt.Sprintf("%v?%v", trafficBaseTongChenYiLong, v.Encode())

	data, err := utils.Request(URL)
	if err != nil {
		return
	}

	var response model.RspTongChengYiLongTrainSeatPrice
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	if len(response.Data.Trains) == 0 {
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}

	var list = make([]model.TrainStationSeatPrice, 0)
	for _, item := range response.Data.Trains {
		for _, seats := range item.TicketState {
			price, err := strconv.ParseFloat(seats.Price, 64)
			if err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			list = append(list, model.TrainStationSeatPrice{
				Train:             item.TrainNum,
				Date:              info.Date,
				StartStation:      info.From,
				StartStationCode:  redis.Get(info.From),
				ArriveStation:     info.To,
				ArriveStationCode: redis.Get(info.To),
				SeatCategoryName:  seats.Cn,
				Price:             int(price * 100),
			})
		}
	}

	cityUpdate(info.From, response.Data.FromCityName)
	cityUpdate(info.To, response.Data.ToCityName)
	seatPriceChan <- SeatPriceInfo{Date: info.Date, From: info.From, To: info.To, Beans: list}
}

func ZhiXingConsumer() {
	ticker := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-ticker:
			num = 0
		default:
			if num > 10 {
				continue
			}
			info := <-seatRelationChan
			go zhiXingTraffic(info)
			num += 1
		}
	}
}

func zhiXingTraffic(info SeatRelationInfo) {
	trafficBaseZhiXing := "http://m.suanya.com/restapi/soa2/10103/json/GetBookingByStationV3ForPC"

	payload := strings.NewReader(fmt.Sprintf("{\"DepartStation\":\"%v\",\"ArriveStation\":\"%v\",\"DepartDate\":\"%v\"}",
		info.From, info.To, info.Date))

	data, err := utils.PostRequest("POST", trafficBaseZhiXing, payload)
	if err != nil {
		return
	}

	var response model.RspZhiXingTrainSeatPrice
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}
	if len(response.ResponseBody.TrainItems) == 0 {
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}

	var list = make([]model.TrainStationSeatPrice, 0)
	for _, item := range response.ResponseBody.TrainItems {
		for _, seats := range item.TicketResult.TicketItems {
			list = append(list, model.TrainStationSeatPrice{
				Train:             item.TrainName,
				Date:              info.Date,
				StartStation:      info.From,
				StartStationCode:  redis.Get(info.From),
				ArriveStation:     info.To,
				ArriveStationCode: redis.Get(info.To),
				SeatCategoryName:  seats.SeatTypeName,
				Price:             int(seats.Price * 100),
			})
		}
	}

	cityUpdate(info.From, response.ResponseBody.DepartureCity.CityName)
	cityUpdate(info.To, response.ResponseBody.ArriveCity.CityName)
	seatPriceChan <- SeatPriceInfo{Date: info.Date, From: info.From, To: info.To, Beans: list}
}

func QunarConsumer() {
	ticker := time.Tick(time.Second)
	var num = 0
	for {
		select {
		case <-ticker:
			num = 0
		default:
			if num > 0 {
				continue
			}
			info := <-seatRelationChan
			go qunarTraffic(info)
			num += 1
		}
	}
}

func qunarTraffic(info SeatRelationInfo) {
	URL := fmt.Sprintf("https://train.qunar.com/dict/open/s2s.do?dptStation=%v&arrStation=%v&date=%v&user=neibu",
		info.From, info.To, info.Date)
	data, err := utils.Request(URL)
	if err != nil {
		return
	}
	var response model.RspQunarTrainSeatPrice
	err = jsoniter.Unmarshal(data, &response)
	if err != nil {
		log.Logger.Warn(err.Error())
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}

	if len(response.Data.S2SBeanList) == 0 {
		redis.LPush(seatConsumerKey, info.ToString())
		return
	}

	var list = make([]model.TrainStationSeatPrice, 0)
	for _, seats := range response.Data.S2SBeanList {
		for key, value := range seats.Seats {
			var seat = model.TrainStationSeatPrice{
				Train:             seats.TrainNo,
				Date:              info.Date,
				StartStation:      seats.DptStationName,
				StartStationCode:  seats.DptStationCode,
				ArriveStation:     seats.ArrStationName,
				ArriveStationCode: seats.ArrStationCode,
				SeatCategoryName:  key,
				Price:             int(value.Price * 100),
			}
			list = append(list, seat)
		}
	}

	cityUpdate(info.From, response.Data.DptCityName)
	cityUpdate(info.To, response.Data.ArrCityName)

	seatPriceChan <- SeatPriceInfo{Date: info.Date, From: info.From, To: info.To, Beans: list}
}

func TrainSeatPrice(date string) {
	var relations = make([]model.TrainStationRelation, 0)
	if err := db.DB.Where("date = ?", date).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var relationMap = make(map[string]map[int]model.TrainStationRelation)
	for _, relation := range relations {
		if _, exist := relationMap[relation.Train]; !exist {
			relationMap[relation.Train] = make(map[int]model.TrainStationRelation)
		}
		relationMap[relation.Train][relation.StationNumber] = relation
	}

	var mp = make(map[SeatRelationInfo]bool)
	for _, from := range relations {
		if from.StationNumber != 1 && from.ArriveTime == from.StartTime {
			continue
		}
		for number, to := range relationMap[from.Train] {
			if number <= from.StationNumber {
				continue
			}
			var info = SeatRelationInfo{
				Train:             from.Train,
				Date:              date,
				From:              from.Station,
				FromStationNumber: from.StationNumber,
				To:                to.Station,
				ToStationNumber:   to.StationNumber,
			}
			mp[info] = false
		}
	}

	var fromToWayMap = make(map[string]map[string]map[string]SeatRelationInfo)
	for info := range mp {
		if _, exist := fromToWayMap[info.From]; !exist {
			fromToWayMap[info.From] = make(map[string]map[string]SeatRelationInfo)
		}
		if _, ok := fromToWayMap[info.From][info.To]; !ok {
			fromToWayMap[info.From][info.To] = make(map[string]SeatRelationInfo)
		}
		fromToWayMap[info.From][info.To][info.Train] = info
	}

	trainSeatTriggerKey := fmt.Sprintf("%v:%v", config.TrainSeatTriggerPrefix, date)
	for from := range fromToWayMap {
		for to := range fromToWayMap[from] {
			key := fmt.Sprintf("%v:%v", from, to)
			if redis.HashExist(trainSeatTriggerKey, key) {
				continue
			}
			info := SeatRelationInfo{Date: date, From: from, To: to}
			redis.LPush(seatConsumerKey, info.ToString())
			redis.HashAdd(trainSeatTriggerKey, key, time.Now().String())
		}
	}
}

func SeatConsumer() {
	for {
		result := redis.BRPop(seatConsumerKey)
		if len(result) == 0 {
			continue
		}
		var info = SeatRelationInfo{}
		err := jsoniter.Unmarshal([]byte(result[1]), &info)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		seatRelationChan <- info
	}

}

func cityUpdate(name, city string) {
	if name == "" || city == "" {
		return
	}
	if redis.HashExist(config.TrainStationCityUpdateKey, name) {
		return
	}
	sql := fmt.Sprintf("update railway_station set city='%v' where name = '%v'", city, name)
	err := db.DB.Exec(sql).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
	redis.HashAdd(config.TrainStationCityUpdateKey, name, city)
}

func SeatInsertConsumer() {
	tick := time.Tick(5 * time.Second)
	var list = make([]SeatPriceInfo, 0)
	for {
		select {
		case <-tick:
			var beans = make([]model.TrainStationSeatPrice, 0)
			for _, info := range list {
				beans = append(beans, info.Beans...)
			}
			batchInsertSeat(beans)
			list = make([]SeatPriceInfo, 0)
		default:
			info := <-seatPriceChan
			list = append(list, info)
		}
	}
}

func batchInsertSeat(list []model.TrainStationSeatPrice) {
	beans := make([][]interface{}, 0)

	for _, item := range list {
		beans = append(beans, []interface{}{item.Train, item.Date, item.StartStation, item.StartStationCode, item.ArriveStation,
			item.ArriveStationCode, item.SeatCategoryName, item.Price})
	}

	var uniqueFields = []string{"train", "date", "start_station", "arrive_station", "seat_category_name"}
	var updateFields = []string{"price"}
	var allFields = []string{"train", "date", "start_station", "start_station_code", "arrive_station", "arrive_station_code", "seat_category_name", "price"}
	count, err := db.Load((&model.TrainStationSeatPrice{}).TableName(), allFields, uniqueFields, updateFields, beans)
	if err != nil {
		log.Logger.Error("load to db error")
	} else {
		log.Logger.Info("load to db", zap.Int("records", count))
	}
	return
}
