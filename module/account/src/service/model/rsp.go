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

type RspUserData struct {
	UserId    string `json:"user_id"`
	LoginName string `json:"login_name"`
	Username  string `json:"username"`
	SchoolId  string `json:"school_id"`
}

type RspUserDetail struct {
	UserId     string `json:"user_id"`
	Username   string `json:"username"`
	Gender     int    `json:"gender"`
	Age        int    `json:"age"`
	HeadImgUrl string `json:"head_img_url"`
	Phone      string `json:"phone"`
	Verify     int    `json:"verify"`
	SchoolId   string `json:"school_id"`
}

type RspLogin struct {
	RspSuccess
	Data LoginResult `json:"data"`
}

type LoginResult struct {
	UserInfo RspUserData         `json:"user_info"`
	Token    string              `json:"token"`
	Policies map[string][]string `json:"policies"`
}

type SchoolDepartments struct {
	DepartmentId     string              `json:"department_id"`
	DepartmentName   string              `json:"department_name"`
	ChildDepartments []SchoolDepartments `json:"child_departments"`
}

type RspSchoolDepartment struct {
	SchoolId    string              `json:"school_id"`
	SchoolName  string              `json:"school_name"`
	Departments []SchoolDepartments `json:"departments"`
}
