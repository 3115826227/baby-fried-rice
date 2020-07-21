package model

type RespOrganize struct {
	Code int64 `json:"code"`
	Data []struct {
		Children   []interface{} `json:"children"`
		Count      int64         `json:"count"`
		CreateTime string        `json:"create_time"`
		ID         string        `json:"id"`
		Label      string        `json:"label"`
		ParentID   string        `json:"parent_id"`
		SchoolID   string        `json:"school_id"`
		Status     bool          `json:"status"`
		UpdateTime string        `json:"update_time"`
	} `json:"data"`
	Message string `json:"message"`
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
