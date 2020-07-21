package model

type ReqOfficialGroupAdd struct {
	Organize string `json:"organize" binding:"required"`
	Admin    string `json:"admin"`
	Name     string `json:"name" binding:"required"`
}
