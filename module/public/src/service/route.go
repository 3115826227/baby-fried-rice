package service

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/job"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/query"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	group := engine.Group("/api/public")

	group.Use(middleware.MiddlewareSetUserMeta())

	job.Register(group)
	query.Register(group)
}
