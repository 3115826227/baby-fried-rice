package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"net/http"
	"fmt"
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
select subject.* from subject inner join course
on subject.id = course.subject_id where course.grade_id = %v`, gradeId)
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

	SuccessResp(c, "", result)
}
