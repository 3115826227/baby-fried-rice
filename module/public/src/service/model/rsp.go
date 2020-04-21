package model

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspGrade struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspSubject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspSubjects []RspSubject

func (rsp RspSubjects) Len() int {
	return len(rsp)
}

func (rsp RspSubjects) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspSubjects) Less(i, j int) bool {
	return rsp[i].Id < rsp[j].Id
}

type RspCourse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspTutors struct {
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
	Data     []RspTutor `json:"data"`
}

type RspTutor struct {
	Id       int    `json:"id"`
	Title    string `json:"name"`
	Salary   int    `json:"salary"`
	Describe string `json:"describe"`
	Subject  string `json:"subject"`
	Grade    string `json:"grade"`
	Area     string `json:"area"`
	UserId   string `json:"user_id"`
}

type RspSalary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RspSalaries []RspSalary

func (rsp RspSalaries) Len() int {
	return len(rsp)
}

func (rsp RspSalaries) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspSalaries) Less(i, j int) bool {
	return rsp[i].Id < rsp[j].Id
}

type RspStreet struct {
	Street string `json:"street"`
}

type RspLocal struct {
	Local   string      `json:"local"`
	Code    string      `json:"code"`
	Streets []RspStreet `json:"streets"`
}

type RspCity struct {
	City   string     `json:"city"`
	Code   string     `json:"code"`
	Locals []RspLocal `json:"locals"`
}

type RspArea struct {
	Province string    `json:"province"`
	Code     string    `json:"code"`
	Cities   []RspCity `json:"cities"`
}

type RspTrainMeta struct {
	//车次
	Train string `json:"train"`
	//日期
	Date string `json:"date"`
	//途径信息
	Detail []RspTrainDetail `json:"detail"`
	//版本
	Version float32 `json:"version"`
}

type RspTrainDetail struct {
	Station    string `json:"station"`
	Number     int    `json:"number"`
	ArriveTime string `json:"arrive_time"`
	StartTime  string `json:"start_time"`
	StopMinute int    `json:"stop_minute"`
	OverDay    int    `json:"over_day"`
}

type RspStationCity struct {
	Name string `json:"name"`
}

type RspStation struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type RspStationDetail struct {
	City     string       `json:"city"`
	Stations []RspStation `json:"stations"`
}

type RspTrainSeatPrice struct {
	SeatCategory string  `json:"seat_category"`
	Price        float32 `json:"price"`
}

type RspTrainSeatDetail struct {
	Train                string              `json:"train"`
	Date                 string              `json:"date"`
	StartStation         string              `json:"start_station"`
	StartTime            string              `json:"start_time"`
	ArriveStation        string              `json:"arrive_station"`
	ArriveTime           string              `json:"arrive_time"`
	FromStation          string              `json:"from_station"`
	FromStationStart     string              `json:"from_station_start"`
	FromStationStartDate string              `json:"from_station_start_date"`
	FromStationNumber    int                 `json:"from_station_number"`
	ToStation            string              `json:"to_station"`
	ToStationArrive      string              `json:"to_station_arrive"`
	ToStationArriveDate  string              `json:"to_station_arrive_date"`
	ToStationNumber      int                 `json:"to_station_number"`
	RunningMinute        int                 `json:"running_minute"`
	Seats                []RspTrainSeatPrice `json:"seats"`
}

type RspTrainSeat struct {
	QueryDate  string               `json:"query_date"`
	FromDetail RspStationDetail     `json:"from_detail"`
	ToDetail   RspStationDetail     `json:"to_detail"`
	Trains     []RspTrainSeatDetail `json:"trains"`
}
