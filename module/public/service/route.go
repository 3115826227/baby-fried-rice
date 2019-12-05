package service

import (
	"github.com/3115826227/baby-fried-rice/module/public/service/job"
	"github.com/3115826227/baby-fried-rice/module/public/service/query"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	group := engine.Group("/api/public")

	job.Register(group)
	query.Register(group)
}
