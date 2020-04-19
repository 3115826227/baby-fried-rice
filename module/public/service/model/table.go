package model

import (
	"encoding/json"
	"time"
)

type CommonIntField struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type CommonStringField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type MessageType struct {
	ID   int `gorm:"AUTO_INCREMENT;column:id;"`
	Name string
}

func (messageType *MessageType) TableName() string {
	return "public_job_message_type"
}

type Subject struct {
	ID   int    `gorm:"AUTO_INCREMENT;column:id;"`
	Name string `gorm:"column:name;unique"`
}

func (table *Subject) TableName() string {
	return "public_job_tutor_subject"
}

type Grade struct {
	ID   int    `gorm:"AUTO_INCREMENT;column:id;"`
	Name string `gorm:"column:name;unique"`
}

func (table *Grade) TableName() string {
	return "public_job_tutor_grade"
}

type Grades []Grade

func (grades Grades) Len() int {
	return len(grades)
}

func (grades Grades) Swap(i, j int) {
	grades[i], grades[j] = grades[j], grades[i]
}

func (grades Grades) Less(i, j int) bool {
	return grades[i].ID < grades[j].ID
}

type Course struct {
	ID        int `gorm:"AUTO_INCREMENT;column:id;"`
	SubjectId int `gorm:"unique(one_record)"`
	GradeId   int `gorm:"unique(one_record)"`
	Name      string
}

func (table *Course) TableName() string {
	return "public_job_tutor_course"
}

//家教
type Tutor struct {
	CommonIntField

	//家教标题
	Title string
	//家教薪资 / hour
	Salary int
	//具体描述
	Describe string
	//科目
	SubjectId int
	//年级
	GradeId int
	//地区
	Area string
	//创建id
	UserId string
	//信息状态
	Status int
	//紧急状态
	Emergency bool
}

func (tutor *Tutor) TableName() string {
	return "public_job_tutor"
}

type Appointment struct {
	CommonStringField

	MessageTypeID int    `xorm:"column:message_type_id;type:int;unique(one_record)"`
	MessageID     string `xorm:"column:message_id;type:char(36);unique(one_record)"`
	AppointmentID string `xorm:"column:appointment_id;type:char(36;)"`
}

func (appointment *Appointment) TableName() string {
	return "public_job_tutor_appointment"
}

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}

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
