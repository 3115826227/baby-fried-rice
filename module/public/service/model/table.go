package model

import "time"

type CommonIntField struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp with time zone" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp with time zone" json:"-"`
}

type CommonStringField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp with time zone" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp with time zone" json:"-"`
}

type MessageType struct {
	ID   int `gorm:"AUTO_INCREMENT;column:id;"`
	Name string
}

func (messageType *MessageType) TableName() string {
	return "job_message_type"
}

type Subject struct {
	ID   int    `gorm:"AUTO_INCREMENT;column:id;"`
	Name string `gorm:"column:name;unique"`
}

type Grade struct {
	ID   int    `gorm:"AUTO_INCREMENT;column:id;"`
	Name string `gorm:"column:name;unique"`
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
	AreaId string
	//创建id
	UserId string
	//信息状态
	Status int
}

func (tutor *Tutor) TableName() string {
	return "job_tutor"
}

type Appointment struct {
	CommonStringField

	MessageTypeID int    `xorm:"column:message_type_id;type:int;unique(one_record)"`
	MessageID     string `xorm:"column:message_id;type:char(36);unique(one_record)"`
	AppointmentID string `xorm:"column:appointment_id;type:char(36;)"`
}

func (appointment *Appointment) TableName() string {
	return "job_appointment"
}

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}

type Station struct {
	Id   int
	Name string
	Code string
	City string
}

type StationCity struct {
	City string
}

type StationCities []StationCity

func (cities StationCities) Len() int {
	return len(cities)
}

func (cities StationCities) Swap(i, j int) {
	cities[i], cities[j] = cities[j], cities[i]
}

func (cities StationCities) Less(i, j int) bool {
	return cities[i].City < cities[j].City
}

type TrainStationMetaRelation struct {
	Train         string
	Station       string
	StationCode   string
	StationNumber int
	ArriveTime    string
	StartTime     string
	StopMinute    int
	OverDay       int
}

type TrainMeta struct {
	Train             string `json:"train"`
	TrainCode         string `json:"train_code"`
	StartStation      string `json:"start_station"`
	StartStationCode  string `json:"start_station_code"`
	ArriveStation     string `json:"arrive_station"`
	ArriveStationCode string `json:"arrive_station_code"`
	StartTime         string `json:"start_time"`
	ArriveTime        string `json:"arrive_time"`
	RunningMinute     int    `json:"running_minute"`
	OverDay           int    `json:"over_day"`
}
