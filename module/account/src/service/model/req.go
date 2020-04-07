package model

type ReqPolicyAdd struct {
}

type ReqRoleAdd struct {
	Name     string   `json:"name" binding:"required"`
	Describe string   `json:"describe"`
	Policies []string `json:"policies"`
}

type ReqLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
}

type ReqAdminLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
}

type ReqAdminAdd struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
}

type ReqAdminInit struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
}

type ReqClientAdd struct {
	Name      string   `json:"name" binding:"required"`
	SchoolIds []string `json:"school_id" binding:"required"`
}

type ReqUserRegister struct {
	ReqLogin
	Username string `json:"username" binding:"required"` //昵称
	Gender   bool   `json:"gender" binding:"required"`   //性别
}

type ReqUserAdd struct {
	ReqLogin

	Roles    []string `json:"roles"`
	SchoolId string   `json:"school_id" binding:"required"`
}

type ReqUserUpdate struct {
	Username string `json:"username"`
}

type ReqUserVerify struct {
	SchoolId string `json:"school_id"`
	Identify string `json:"identify"`
	Name     string `json:"name"`
}

type ReqSchoolDepartmentAdd struct {
	SchoolId string `json:"school_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	ParentId string `json:"parent_id"`
}

type ReqSchoolDepartmentUpdate struct {
	SchoolDepartmentId string `json:"school_department_id" binding:"required"`
	ReqSchoolDepartmentAdd
}

type ReqSchoolCertificationDelete struct {
	Id []string `json:"id"`
}

type ReqSchoolOrganizeAdd struct {
	Label    string `json:"label" binding:"required"`
	ParentId string `json:"parent_id" binding:"required"`
	SchoolId string `json:"school_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

type ReqSchoolOrganizedUpdate struct {
	Id    string `json:"id" binding:"required"`
	Label string `json:"label" binding:"required"`
}

type ReqSchoolOrganizedStatusUpdate struct {
	Id     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type ReqSchoolStudentAdd struct {
	Organize string `json:"organize"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Identify string `json:"identify"`
}
