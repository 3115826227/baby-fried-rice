package model

import (
	"time"
	"encoding/json"
)

type CommonField struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp with time zone" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp with time zone" json:"-"`
}

//火车站信息
type Station struct {
	CommonField
	Name string `gorm:"unique_index:idx_station_name_code_name_code"`
	Code string `gorm:"unique_index:idx_station_name_code_name_code"`
	City string
}

func (station *Station) TableName() string {
	return "railway_station"
}

type Stations []Station

func (stations Stations) Len() int {
	return len(stations)
}

func (stations Stations) Swap(i, j int) {
	stations[i], stations[j] = stations[j], stations[i]
}

func (stations Stations) Less(i, j int) bool {
	return stations[i].ID < stations[j].ID
}

type TrainMeta struct {
	Train                string `gorm:"unique_index:idx_train_date_train_date"`
	TrainCode            string
	Date                 string `gorm:"unique_index:idx_train_date_train_date"`
	StartStation         string
	StartStationCode     string
	ArriveStation        string
	ArriveStationCode    string
	StartTime            string
	ArriveTime           string
	RunningStationNumber int
	RunningMinute        int
	OverDay              int
}

func (trainMeta *TrainMeta) TableName() string {
	return "railway_train_meta"
}

type TrainStationRelation struct {
	Train         string `gorm:"column:train;unique_index:idx_train_station_train_station"`
	Date          string `gorm:"column:date;unique_index:idx_train_station_train_station"`
	Station       string
	StationCode   string
	StationNumber int    `gorm:"column:station_number;unique_index:idx_train_station_train_station"`
	ArriveTime    string `gorm:"column:arrive_time;type:char(20)"`
	StartTime     string `gorm:"column:start_time;type:char(20)"`
	StopMinute    int
	OverDay       int
}

func (relation *TrainStationRelation) TableName() string {
	return "railway_ref_train_station"
}

type TrainStationSeatPrice struct {
	Train             string `gorm:"column:train;unique_index:idx_train_station_seat_price"`
	Date              string `gorm:"column:date;unique_index:idx_train_station_seat_price"`
	StartStation      string `gorm:"column:start_station;unique_index:idx_train_station_seat_price"`
	StartStationCode  string
	ArriveStation     string `gorm:"column:arrive_station;unique_index:idx_train_station_seat_price"`
	ArriveStationCode string
	TrainCategoryName string
	SeatCategoryName  string `gorm:"column:seat_category_name;unique_index:idx_train_station_seat_price"`
	Price             int
	StudentPrice      int
	ChildPrice        int
}

func (price *TrainStationSeatPrice) TableName() string {
	return "railway_train_station_seat_price"
}

func (price *TrainStationSeatPrice) ToString() string {
	data, _ := json.Marshal(price)
	return string(data)
}

//列车类型
type TrainCategory struct {
	CommonField
	Name string `gorm:"unique"`
}

func (trainCategory *TrainCategory) TableName() string {
	return "railway_train_category"
}

//座位类型
type TrainSeatCategory struct {
	CommonField
	Name string `gorm:"unique"`
}

func (trainSeatCategory *TrainSeatCategory) TableName() string {
	return "railway_train_seat_category"
}
