package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"github.com/gin-gonic/gin"
)

func TutorGet(c *gin.Context) {
	result := make([]model.RspTutor, 0)
	//page := c.Query("page")
	//pageSize := c.Query("page_size")
	search := c.Query("search")
	var tutors = make([]model.Tutor, 0)
	sql := `select * from public_job_tutor where id like '%` + search + `%'`
	if err := db.DB.Raw(sql).Scan(&tutors).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var subjects = make([]model.Subject, 0)
	var grades = make([]model.Grade, 0)
	var subjectMap = make(map[int]model.Subject)
	var gradeMap = make(map[int]model.Grade)
	if err := db.DB.Find(&subjects).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	if err := db.DB.Find(&grades).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	for _, s := range subjects {
		subjectMap[s.ID] = s
	}
	for _, g := range grades {
		gradeMap[g.ID] = g
	}
	for _, t := range tutors {
		result = append(result, model.RspTutor{
			Id:       t.ID,
			Title:    t.Title,
			Salary:   t.Salary,
			Describe: t.Describe,
			Subject:  subjectMap[t.SubjectId].Name,
			Grade:    gradeMap[t.GradeId].Name,
			Area:     t.Area,
			UserId:   "",
		})
	}
	//if err := db.DB.Where("id LIKE ?", "%"+ search + "%").Find(&tutors).Error;err != nil {
	//	log.Logger.Warn(err.Error())
	//}

	SuccessResp(c, "", result)
}
