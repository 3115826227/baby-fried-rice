package service

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service/job"
	"github.com/3115826227/baby-fried-rice/module/public/service/query"
)

func Register(engine *gin.Engine) {

	group := engine.Group("/api/public")

	job.Register(group)
	query.Register(group)
}
