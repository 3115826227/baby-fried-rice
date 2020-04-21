package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func GradeGet(c *gin.Context) {
	var result = make([]model.RspGrade, 0)
	var grades = make([]model.Grade, 0)

	if err := db.DB.Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	sort.Sort(model.Grades(grades))

	for _, grade := range grades {
		result = append(result, model.RspGrade{Id: grade.ID, Name: grade.Name})
	}

	SuccessResp(c, "", result)
}

func SubjectGet(c *gin.Context) {
	gradeId := c.Query("grade_id")

	var result = make([]model.RspSubject, 0)
	if gradeId == "" {
		var subjects = make([]model.Subject, 0)
		if err := db.DB.Find(&subjects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
		for _, s := range subjects {
			result = append(result, model.RspSubject{Id: s.ID, Name: s.Name})
		}
	} else {
		var sql = fmt.Sprintf(`
select public_job_tutor_subject.* from public_job_tutor_subject inner join public_job_tutor_course
on public_job_tutor_subject.id = public_job_tutor_course.subject_id where public_job_tutor_course.grade_id = %v`, gradeId)
		rows, err := db.DB.Raw(sql).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var rspSubject = new(model.RspSubject)
			err := rows.Scan(&rspSubject.Id, &rspSubject.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, sysErrResponse)
				return
			}
			result = append(result, *rspSubject)
		}
	}

	sort.Sort(model.RspSubjects(result))

	SuccessResp(c, "", result)
}

func CourseGet(c *gin.Context) {
	var result = make([]model.RspCourse, 0)
	var courses = make([]model.Course, 0)
	if err := db.DB.Find(&courses).Error;err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	for _, c := range courses {
		result = append(result, model.RspCourse{
			Id:   c.ID,
			Name: c.Name,
		})
	}

	SuccessResp(c, "", result)
}
