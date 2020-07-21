package model

type RespUserLogin struct {
	Code int64 `json:"code"`
	Data struct {
		Detail struct {
			Addr      string `json:"addr"`
			Age       int64  `json:"age"`
			Birthday  string `json:"birthday"`
			CreatedAt string `json:"created_at"`
			Ethnic    string `json:"ethnic"`
			Gr        bool   `json:"gr"`
			Hometown  string `json:"hometown"`
			ID        string `json:"id"`
			Phone     string `json:"phone"`
			Qq        string `json:"qq"`
			SchoolID  string `json:"school_id"`
			UpdatedAt string `json:"updated_at"`
			Username  string `json:"username"`
			AccountId string `json:"account_id"`
			Verify    bool   `json:"verify"`
			Wx        string `json:"wx"`
		} `json:"detail"`
		User struct {
			CreatedAt  string `json:"created_at"`
			EncodeType string `json:"encode_type"`
			ID         string `json:"id"`
			LoginName  string `json:"login_name"`
			Password   string `json:"password"`
			UpdatedAt  string `json:"updated_at"`
		} `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespUserVerifyDetail struct {
	Code int64 `json:"code"`
	Data struct {
		Detail struct {
			Addr      string `json:"addr"`
			Age       int64  `json:"age"`
			Birthday  string `json:"birthday"`
			CreatedAt string `json:"created_at"`
			Ethnic    string `json:"ethnic"`
			Gender    bool   `json:"gender"`
			Hometown  string `json:"hometown"`
			ID        string `json:"id"`
			Phone     string `json:"phone"`
			Qq        string `json:"qq"`
			SchoolID  string `json:"school_id"`
			UpdatedAt string `json:"updated_at"`
			Username  string `json:"username"`
			AccountId string `json:"account_id"`
			Verify    bool   `json:"verify"`
			Wx        string `json:"wx"`
		} `json:"detail"`
		School struct {
			CreatedAt string `json:"created_at"`
			ID        string `json:"id"`
			Identify  string `json:"identify"`
			Name      string `json:"name"`
			Number    string `json:"number"`
			OrgID     string `json:"org_id"`
			UpdatedAt string `json:"updated_at"`
		} `json:"school"`
		SchoolLabel struct {
			Faculty string `json:"faculty"`
			Grade   string `json:"grade"`
			Major   string `json:"major"`
			School  string `json:"school"`
		} `json:"school_label"`
		User struct {
			CreatedAt  string `json:"created_at"`
			EncodeType string `json:"encode_type"`
			ID         string `json:"id"`
			LoginName  string `json:"login_name"`
			Password   string `json:"password"`
			UpdatedAt  string `json:"updated_at"`
		} `json:"user"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespUserInfo struct {
	Code int64 `json:"code"`
	Data struct {
		Age       int64  `json:"age"`
		Faculty   string `json:"faculty"`
		Gender    bool   `json:"gender"`
		Grade     string `json:"grade"`
		ID        string `json:"id"`
		Major     string `json:"major"`
		Name      string `json:"name"`
		Number    string `json:"number"`
		School    string `json:"school"`
		Username  string `json:"username"`
		AccountId string `json:"account_id"`
		Verify    bool   `json:"verify"`
	} `json:"data"`
	Message string `json:"message"`
}

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
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
	SchoolId  string `json:"school_id"`
	IsSuper   bool   `json:"is_super"`
}

type RspUserDetail struct {
	UserId    string `json:"user_id"`
	AccountId string `json:"account_id"`
	LoginName string `json:"login_name"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Gender    bool   `json:"gender"`
	Age       int64  `json:"age"`
	//HeadImgUrl string `json:"head_img_url"`
	Phone    string `json:"phone"`
	Verify   bool   `json:"verify"`
	School   string `json:"school"`
	Faculty  string `json:"faculty"`
	Grade    string `json:"grade"`
	Major    string `json:"major"`
	Number   string `json:"number"`
	Identify string `json:"identify"`
	SchoolId string `json:"school_id"`
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

type RspAdminPermission struct {
	Id       int                  `json:"id"`
	Name     string               `json:"name"`
	Method   string               `json:"method"`
	Path     string               `json:"path"`
	Types    int                  `json:"types"`
	Children []RspAdminPermission `json:"children"`
	ParentId int                  `json:"-"`
}

type RspAdminPermissions []RspAdminPermission

func (rsp RspAdminPermissions) Len() int {
	return len(rsp)
}

func (rsp RspAdminPermissions) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspAdminPermissions) Less(i, j int) bool {
	return rsp[i].Id < rsp[j].Id
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

type RspSchoolOrganizeStudent struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
	Label      string `json:"label"`
	Identify   string `json:"identify"`
	Phone      string `json:"phone"`
	Status     string `json:"verify"`
	UpdateTime string `json:"update_time"`
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
