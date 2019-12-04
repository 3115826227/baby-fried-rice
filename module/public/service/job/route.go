package job

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service/job/handle"
)

func Register(engine *gin.RouterGroup) {

	engine.POST("/tutor", handle.TutorAdd)

	engine.GET("/grade", handle.GradeGet)
	engine.GET("/subject", handle.SubjectGet)
}
