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
	IsSuper   bool   `json:"is_super"`
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
	UserInfo   RspUserData `json:"user_info"`
	Token      string      `json:"token"`
	Role       []AdminRole `json:"role"`
	Permission []int       `json:"permission"`
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

type RspSubAdmin struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type RspAdminPermissions struct {
	Id       int                   `json:"id"`
	Name     string                `json:"name"`
	Method   string                `json:"method"`
	Path     string                `json:"path"`
	Types    int                   `json:"types"`
	Children []RspAdminPermissions `json:"children"`
	ParentId int                   `json:"-"`
}

type RspImages struct {
	Id         string   `json:"id"`
	Name       []string `json:"name"`
	Size       int64    `json:"size"`
	Timestamp  int64    `json:"timestamp"`
	CreateTime string   `json:"create_time"`
}

type RspContainers struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	ImageId    string           `json:"image_id"`
	ImageName  string           `json:"image_name"`
	Timestamp  int64            `json:"timestamp"`
	CreateTime string           `json:"create_time"`
	State      string           `json:"state"`
	Status     string           `json:"status"`
	Ports      []ContainerPorts `json:"ports"`
}

type ContainerPorts struct {
	Ip          string `json:"ip"`
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port"`
}

type RspCmdContainerStats struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	UseLimitMemory string `json:"use_limit_memory"`
	MemoryPercent  string `json:"memory_percent"`
	CpuPercent     string `json:"cpu_percent"`
}

type RspContainerStats struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	UseMemory        uint64 `json:"use_memory"`
	LimitMemory      uint64 `json:"limit_memory"`
	UseMemoryPercent string `json:"use_memory_percent"`
	UseCPUPercent    string `json:"use_cpu_percent"`
}
