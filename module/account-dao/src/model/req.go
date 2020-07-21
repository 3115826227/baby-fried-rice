package model

type ReqPasswordLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type ReqPhoneLogin struct {
	Phone string `json:"phone" binding:"required"`
}

type ReqRoleAdd struct {
	Name        string `json:"name" binding:"required"`
	SchoolId    string `json:"school_id" binding:"required"`
	Describe    string `json:"describe"`
	Permissions []int  `json:"permissions" binding:"required"`
}

type ReqRoleUpPermission struct {
	Role       int  `json:"role" binding:"required"`
	Permission int  `json:"permission" binding:"required"`
	Status     bool `json:"status" binding:"required"`
}

type ReqAdminInit struct {
	LoginName string `json:"login_name" binding:"required"` // 账号登陆名称
	Name      string `json:"name" binding:"required"`       //账号名称
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
}

type ReqPasswordAdminLogin struct {
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type ReqSubAdminAdd struct {
	ReqId     string `json:"req_id" binding:"required"`
	LoginName string `json:"login_name" binding:"required"` // 账号名称
}

type ReqSubAdminUpRole struct {
	Admin  string `json:"admin" binding:"required"`
	Role   int    `json:"role" binding:"required"`
	Status bool   `json:"status" binding:"required"`
}

type ReqSchoolStudentAdd struct {
	Organize string `json:"organize"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Identify string `json:"identify"`
}

type ReqSchoolOrganizeAdd struct {
	Label    string `json:"label" binding:"required"`
	ParentId string `json:"parent_id" binding:"required"`
	SchoolId string `json:"school_id" binding:"required"`
	Status   bool   `json:"status" binding:"required"`
}

type ReqSchoolOrganizedUpdate struct {
	Id    string `json:"id" binding:"required"`
	Label string `json:"label" binding:"required"`
}

type ReqSchoolOrganizedStatusUpdate struct {
	Id     string `json:"id" binding:"required"`
	Status bool   `json:"status" binding:"required"`
}

type ReqSchoolAdd struct {
	Name     string `json:"name"`
	Province string `json:"province"`
	City     string `json:"city"`
}

type ReqUserVerify struct {
	UserId   string `json:"user_id"`
	SchoolId string `json:"school_id"`
	Identify string `json:"identify"`
	Name     string `json:"name"`
}

type ReqUserUpdate struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type ReqUserRegister struct {
	ReqPasswordLogin
	Username string `json:"username" binding:"required"` //昵称
	Gender   bool   `json:"gender" binding:"required"`   //性别
}
