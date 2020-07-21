package model

type ReqSchoolStudentAdd struct {
	Organize string `json:"organize"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Identify string `json:"identify"`
}

type RespSchoolStudent struct {
	Code int64 `json:"code"`
	Data []struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Identify   string `json:"identify"`
		Status     bool   `json:"status"`
		Number     string `json:"number"`
		Phone      string `json:"phone"`
		OrgId      string `json:"org_id"`
		CreateTime string `json:"create_time"`
		UpdateTime string `json:"update_time"`
	} `json:"data"`
	Message string `json:"message"`
}
