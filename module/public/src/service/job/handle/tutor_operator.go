package handle

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/config"
	"github.com/3115826227/baby-fried-rice/module/public/src/log"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func IsValidCourse(name string) (model.Course, bool) {
	var course model.Course
	var count = 0
	if err := db.DB.Where("name = ?", name).First(&course).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return model.Course{}, false
	}
	return course, count != 0
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
	salary, err := strconv.Atoi(req.Salary)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	if req.Title == "" || salary == 0 || req.Area == "" || req.Describe == "" || req.Emergency == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	course, exist := IsValidCourse(req.Course)
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var tutor model.Tutor
	var now = time.Now()
	//tutor.UserId = GetUserID(c)
	tutor.CreatedAt, tutor.UpdatedAt = now, now
	tutor.Title = req.Title
	tutor.SubjectId = course.SubjectId
	tutor.GradeId = course.GradeId
	tutor.Describe = req.Describe
	tutor.Area = req.Area
	tutor.Salary = salary
	if req.Emergency == config.TutorEmergency {
		tutor.Emergency = true
	}

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
