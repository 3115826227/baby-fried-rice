package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/public/src/config"
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func TutorGet(c *gin.Context) {
	result := make([]model.RspTutor, 0)
	page := c.Query("page")
	pageSize := c.Query("page_size")
	grade := c.Query("grade_id")
	subject := c.Query("subject_id")
	salary := c.Query("salary_id")
	search := c.Query("search")
	var pageInt, pageSizeInt int
	var err error
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
	} else {
		pageInt = config.DefaultPage
	}
	if pageSize != "" {
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil || pageSizeInt <= 0 {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
	} else {
		pageSizeInt = config.DefaultPageSize
	}
	if grade != "" {
		_, err = strconv.Atoi(grade)
		if err != nil {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
	}
	if subject != "" {
		_, err = strconv.Atoi(subject)
		if err != nil {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
	}
	if salary != "" {
		_, err = strconv.Atoi(salary)
		if err != nil {
			log.Logger.Warn(err.Error())
			c.JSON(http.StatusBadRequest, paramErrResponse)
			return
		}
	}

	var tutors = make([]model.Tutor, 0)
	var count = struct {
		Count int `json:"count"`
	}{}
	countSql := `select count(1) as count from public_job_tutor as tutor
left join public_job_tutor_course as course
on course.grade_id = tutor.grade_id and course.subject_id = tutor.subject_id
left join public_job_tutor_salary as salary
on tutor.salary > salary.min and tutor.salary <= salary.max
where tutor.id like '%` + search + `%'`
	sql := `select tutor.* from public_job_tutor as tutor
left join public_job_tutor_course as course
on course.grade_id = tutor.grade_id and course.subject_id = tutor.subject_id
left join public_job_tutor_salary as salary
on tutor.salary > salary.min and tutor.salary <= salary.max
where tutor.id like '%` + search + `%'`
	if subject != "" {
		sql += ` and course.subject_id = ` + subject
		countSql += ` and course.subject_id = ` + subject
	}
	if grade != "" {
		sql += ` and course.grade_id = ` + grade
		countSql += ` and course.grade_id = ` + grade
	}
	if salary != "" {
		sql += ` and salary.id = ` + salary
		countSql += ` and salary.id = ` + salary
	}
	if err := db.DB.Debug().Raw(countSql).Scan(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	sql += ` order by create_time desc limit ` +
		fmt.Sprintf("%v", pageInt-1) + `,` + fmt.Sprintf("%v", pageSizeInt)
	if err := db.DB.Debug().Raw(sql).Scan(&tutors).Error; err != nil {
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
			UserId:   t.UserId,
		})
	}

	rsp := model.RspTutors{
		Page:     pageInt,
		PageSize: pageSizeInt,
		Total:    count.Count,
		Data:     result,
	}
	fmt.Println(rsp)

	SuccessResp(c, "", rsp)
}
