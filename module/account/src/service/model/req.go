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

type ReqClientAdd struct {
	Name      string   `json:"name" binding:"required"`
	SchoolIds []string `json:"school_id" binding:"required"`
}

type ReqUserRegister struct {
	ReqLogin
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
