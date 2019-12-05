package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func IsValidGrade(id int) bool {
	var grade model.Grade
	var count = 0
	if err := db.DB.Where("id = ?", id).First(&grade).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return false
	}
	return count != 0
}

func IsValidSubject(id int) bool {
	var subject model.Subject
	var count = 0
	if err := db.DB.Where("id = ?", id).First(&subject).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return false
	}
	return count != 0
}

func IsValidCourse(subjectId, gradeId int) bool {
	var course model.Course
	var count = 0
	if err := db.DB.Where("subject_id = ? and grade_id = ?", subjectId, gradeId).First(&course).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return false
	}
	return count != 0
}

func IsValidArea(code string) bool {
	var area model.Area
	var count = 0
	if err := db.DB.Where("code = ? and level = 3", code).First(&area).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return false
	}
	return count != 0
}

func GetUserID(c *gin.Context) string {
	return c.GetHeader("userId")
}

func TutorAdd(c *gin.Context) {
	var err error
	var req model.ReqTutorAdd
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if req.Title == "" || req.Salary == 0 || req.GradeId == 0 ||
		req.SubjectId == 0 || req.AreaId == "" || req.Describe == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if !IsValidCourse(req.SubjectId, req.GradeId) || !IsValidArea(req.AreaId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var tutor model.Tutor
	var now = time.Now()
	tutor.UserId = GetUserID(c)
	tutor.CreatedAt, tutor.UpdatedAt = now, now
	tutor.Title = req.Title
	tutor.SubjectId = req.SubjectId
	tutor.GradeId = req.GradeId
	tutor.Describe = req.Describe
	tutor.AreaId = req.AreaId
	tutor.Salary = req.Salary

	err = db.DB.Create(&tutor).Error
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func TutorUpdate() {

}

func TutorDelete() {

}

func TutorAppointment() {

}

func TutorCancelAppointment() {

}
