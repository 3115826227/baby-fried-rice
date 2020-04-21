package model

type ReqTutorAdd struct {
	Title     string `json:"title"`
	Salary    string `json:"salary"`
	Course    string `json:"course"`
	Area      string `json:"area"`
	Describe  string `json:"describe"`
	Emergency string `json:"emergency"`
}

type ReqTrainSeatGet struct {
	From     string `json:"from" binding:"required"`
	To       string `json:"to" binding:"required"`
	Date     string `json:"date" binding:"required"`
	HighOnly bool   `json:"high_only"`
}
