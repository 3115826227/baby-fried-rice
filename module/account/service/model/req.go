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
	Name      string   `json:"name;" binding:"required"`
	SchoolIds []string `json:"school_id" binding:"required"`
}

type ReqUserRegister struct {
	ReqLogin
}

type ReqUserAdd struct {
	ReqLogin

	Roles    []string `json:"roles"`
	ClientId string   `json:"client_id" binding:"required"`
	SchoolId string   `json:"school_id" binding:"required"`
}
