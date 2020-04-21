package job

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/service/job/handle"
	"github.com/gin-gonic/gin"
)

func init()  {
	//handle.AddSalary()
	//handle.AddGrade()
	//handle.AddCourse()
}

func Register(engine *gin.RouterGroup) {

	tutor := engine.Group("/job/tutor")

	tutor.POST("", handle.TutorAdd)
	tutor.GET("", handle.TutorGet)

	tutor.GET("/grade", handle.GradeGet)
	tutor.GET("/subject", handle.SubjectGet)
	tutor.GET("/course", handle.CourseGet)
	tutor.GET("/salary", handle.GetSalary)
}
