package model

type ReqAdminLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
	Ip        string `json:"ip"`
}

type RespAdminLogin struct {
	Code int64 `json:"code"`
	Data struct {
		Admin struct {
			CreatedAt  string `json:"created_at"`
			EncodeType string `json:"encode_type"`
			ID         string `json:"id"`
			LoginName  string `json:"login_name"`
			Name       string `json:"name"`
			Password   string `json:"password"`
			Phone      string `json:"phone"`
			ReqID      string `json:"req_id"`
			SchoolID   string `json:"school_id"`
			Super      bool   `json:"super"`
			UpdatedAt  string `json:"updated_at"`
			Username   string `json:"username"`
		} `json:"admin"`
		Permissions []int           `json:"permissions"`
		Roles       []RespAdminRole `json:"roles"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespAdminRole struct {
	CreatedAt string `json:"created_at"`
	Describe  string `json:"describe"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SchoolID  string `json:"school_id"`
	UpdatedAt string `json:"updated_at"`
}

type RspUserData struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
	SchoolId  string `json:"school_id"`
	IsSuper   bool   `json:"is_super"`
}

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RspLogin struct {
	RspSuccess
	Data LoginResult `json:"data"`
}

type LoginResult struct {
	UserInfo   RspUserData     `json:"user_info"`
	Token      string          `json:"token"`
	Role       []RespAdminRole `json:"role"`
	Permission []int           `json:"permission"`
}

type AdminRole struct {
	ID        int64
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	SchoolId  string `json:"-"`
}
