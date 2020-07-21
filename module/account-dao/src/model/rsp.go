package model

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RespAdminLogin struct {
	Admin       AccountAdmin      `json:"admin"`
	Roles       []AdminRole       `json:"roles"`
	Permissions []AdminPermission `json:"permissions"`
}

type RespUserLogin struct {
	User   AccountUser       `json:"user"`
	Detail AccountUserDetail `json:"detail"`
}

type RespSchoolLabel struct {
	School  string `json:"school"`
	Faculty string `json:"faculty"`
	Grade   string `json:"grade"`
	Major   string `json:"major"`
}

type RespUserDetail struct {
	User        AccountUser             `json:"user"`
	Detail      AccountUserDetail       `json:"detail"`
	School      AccountUserSchoolDetail `json:"school"`
	SchoolLabel RespSchoolLabel         `json:"school_label"`
}

type RespIp struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

type RspLoginLog struct {
	LoginName string `json:"login_name"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Count     int    `json:"count"`
	Ip        string `json:"ip"`
	Area      string `json:"area"`
	Time      string `json:"time"`
}

type RspUserLoginLog struct {
	RspLoginLog
	UserId string `json:"user_id"`
}

type RspAdminLoginLog struct {
	RspLoginLog
	AdminId string `json:"admin_id"`
	School  string `json:"school"`
}

type RspRootLoginLog struct {
	RspLoginLog
	RootId string `json:"root_id"`
}

type RspUserInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Gender   bool   `json:"gender"`
	Age      int    `json:"age"`
	//HeadImgUrl string `json:"head_img_url"`
	Verify  bool   `json:"verify"`
	School  string `json:"school"`
	Faculty string `json:"faculty"`
	Grade   string `json:"grade"`
	Major   string `json:"major"`
	Number  string `json:"number"`
}

type RspSchoolOrganize struct {
	Id         string              `json:"id"`
	Label      string              `json:"label"`
	ParentId   string              `json:"parent_id"`
	Count      int                 `json:"count"`
	Status     bool                `json:"status"`
	CreateTime string              `json:"create_time"`
	UpdateTime string              `json:"update_time"`
	SchoolId   string              `json:"school_id"`
	Children   []RspSchoolOrganize `json:"children"`
}

type RspSchoolOrganizes []RspSchoolOrganize

func (rsp RspSchoolOrganizes) Len() int {
	return len(rsp)
}

func (rsp RspSchoolOrganizes) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspSchoolOrganizes) Less(i, j int) bool {
	return rsp[i].Id < rsp[j].Id
}

type RspSchoolStudent struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Identify   string `json:"identify"`
	Status     bool   `json:"status"`
	Number     string `json:"number"`
	Phone      string `json:"phone"`
	OrgId      string `json:"org_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type RspSchoolStudents []RspSchoolStudent

func (rsp RspSchoolStudents) Len() int {
	return len(rsp)
}

func (rsp RspSchoolStudents) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspSchoolStudents) Less(i, j int) bool {
	return rsp[i].Number < rsp[j].Number
}
