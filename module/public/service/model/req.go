package model

type ReqTutorAdd struct {
	Title     string `json:"title"`
	Salary    int    `json:"salary"`
	SubjectId int    `json:"subject_id"`
	GradeId   int    `json:"grade_id"`
	AreaId    string `json:"area_id"`
	Describe  string `json:"describe"`
}

type ReqTrainSeatGet struct {
	From     string `json:"from" binding:"required"`
	To       string `json:"to" binding:"required"`
	Date     string `json:"date" binding:"required"`
	HighOnly bool   `json:"high_only"`
}
