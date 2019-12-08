package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/public/config"
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func TrainMetaGet(c *gin.Context) {
	code := c.Query("code")
	date := c.Query("date")
	if code == "" || date == "" {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	sql := fmt.Sprintf(`select station, station_number,arrive_time, start_time,stop_minute, over_day 
from railway_ref_train_station where train = '%v' and date = '%v' order by station_number`, code, date)

	rows, err := db.DB.Raw(sql).Rows()
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	defer rows.Close()

	var rsp = model.RspTrainMeta{}
	var details = make([]model.RspTrainDetail, 0)
	rsp.Train = code
	rsp.Date = date
	for rows.Next() {
		var detail = model.RspTrainDetail{}
		err := rows.Scan(&detail.Station, &detail.Number, &detail.ArriveTime, &detail.StartTime,
			&detail.StopMinute, &detail.OverDay)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		detail.ArriveTime = strings.TrimSpace(detail.ArriveTime)
		detail.StartTime = strings.TrimSpace(detail.StartTime)
		details = append(details, detail)
	}
	rsp.Detail = details

	SuccessResp(c, "", rsp)
}

func IsHighTrain(train string) bool {
	return train[0] == 'G' || train[0] == 'C' || train[0] == 'D'
}

func TrainSeatPriceGet(c *gin.Context) {
	var err error
	var req = model.ReqTrainSeatGet{}
	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, paramErrResponse)
		return
	}

	sql := fmt.Sprintf(`
select 
	train_meta.train, train_meta.date, train_meta.start_station as dpt_station, train_meta.start_time, train_meta.arrive_station as arrive_station, train_meta.arrive_time as arrive_time,
	from_meta.station as from_station,from_meta.station_code as from_station_code, from_meta.start_time as from_station_start, from_meta.station_number as from_station_number,from_meta.over_day as from_over_day,
	to_meta.station as to_station,to_meta.station_code as to_station_code, to_meta.arrive_time as to_station_arrive, to_meta.station_number as to_station_number, to_meta.over_day as to_over_day,
	price.seat_category_name, price.price
from railway_ref_train_station as from_meta
inner join railway_ref_train_station as to_meta on from_meta.train = to_meta.train and from_meta.date = to_meta.date
inner join railway_station as from_station on from_meta.station = from_station.name
inner join railway_station as to_station on to_meta.station = to_station.name
left join railway_train_station_seat_price as price 
on from_meta.train = price.train and from_meta.date = price.date 
and from_meta.station = price.start_station and to_meta.station = price.arrive_station
inner join railway_train_meta as train_meta
on from_meta.train =  train_meta.train and from_meta.date = train_meta.date
where from_meta.station_number < to_meta.station_number and from_meta.date = '%v' and from_station.city = '%v' and to_station.city = '%v'
order by from_meta.start_time
`, req.Date, req.From, req.To)

	rows, err := db.DB.Raw(sql).Rows()
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	defer rows.Close()

	var rsp = model.RspTrainSeat{}
	rsp.QueryDate = req.Date
	var fromStationMap = make(map[string]string)
	var toStationMap = make(map[string]string)
	var trainMap = make(map[string][]model.RspTrainSeatPrice)
	var trainSeatDetailMap = make(map[string]model.RspTrainSeatDetail)
	for rows.Next() {
		var train, date, dptStation, startTime, arriveStation, arriveTime string
		var fromStation, fromStationCode, fromStationStart string
		var fromStationNumber, fromStationOverDay int
		var toStation, toStationCode, toStationArrive string
		var toStationNumber, toStationOverDay int
		var seatCategory string
		var price int
		var detail model.RspTrainSeatDetail
		err = rows.Scan(&train, &date, &dptStation, &startTime, &arriveStation, &arriveTime,
			&fromStation, &fromStationCode, &fromStationStart, &fromStationNumber, &fromStationOverDay,
			&toStation, &toStationCode, &toStationArrive, &toStationNumber, &toStationOverDay,
			&seatCategory, &price)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		if req.HighOnly && !IsHighTrain(train) {
			continue
		}

		fromStationMap[fromStation] = fromStationCode
		toStationMap[toStation] = toStationCode

		key := fmt.Sprintf("%v:%v:%v:%v", train, date, fromStation, toStation)
		if _, exist := trainMap[key]; !exist {
			trainMap[key] = make([]model.RspTrainSeatPrice, 0)
		}
		infos := trainMap[key]
		infos = append(infos, model.RspTrainSeatPrice{SeatCategory: seatCategory, Price: float32(price) / 100})
		trainMap[key] = infos

		tempDate, _ := time.Parse(config.DayLayout, date)

		detail.Train = train
		detail.Date = tempDate.AddDate(0, 0, -fromStationOverDay).Format(config.DayLayout)
		detail.StartStation = dptStation
		detail.StartTime = startTime
		detail.ArriveStation = arriveStation
		detail.ArriveTime = arriveTime
		detail.FromStation = fromStation
		detail.FromStationStart = strings.TrimSpace(fromStationStart)
		detail.FromStationStartDate = date
		detail.FromStationNumber = fromStationNumber
		detail.ToStation = toStation
		detail.ToStationArrive = strings.TrimSpace(toStationArrive)
		detail.ToStationArriveDate = tempDate.AddDate(0, 0, toStationOverDay-fromStationOverDay).Format(config.DayLayout)
		detail.ToStationNumber = toStationNumber
		trainSeatDetailMap[key] = detail
	}

	var fromStationDetail model.RspStationDetail
	fromStationDetail.City = req.From
	for name, code := range fromStationMap {
		fromStationDetail.Stations = append(fromStationDetail.Stations, model.RspStation{Name: name, Code: code})
	}
	var toStationDetail model.RspStationDetail
	toStationDetail.City = req.To
	for name, code := range toStationMap {
		toStationDetail.Stations = append(toStationDetail.Stations, model.RspStation{Name: name, Code: code})
	}
	rsp.FromDetail = fromStationDetail
	rsp.ToDetail = toStationDetail

	var details = make([]model.RspTrainSeatDetail, 0)
	for key, seats := range trainMap {
		detail := trainSeatDetailMap[key]
		detail.Seats = seats
		details = append(details, detail)
	}
	rsp.Trains = details

	SuccessResp(c, "", rsp)

	return

}
