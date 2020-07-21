package service

import (
	"github.com/3115826227/baby-fried-rice/module/public/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/public/src/service/job"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	group := engine.Group("/api/public")

	group.Use(middleware.MiddlewareSetUserMeta())

	//兼职路由
	job.Register(group)

	//表白墙
	group.POST("/confession/wall")
	group.PATCH("/confession/wall")
	group.GET("/confession/wall")
	group.DELETE("/confession/wall")

	//租房信息
	group.POST("/rent/house")
	group.PATCH("/rent/house")
	group.GET("/rent/house")
	group.DELETE("/rent/house")

	//失物招领
	group.POST("/things/lost")
	group.PATCH("/things/lost")
	group.GET("/things/lost")
	group.POST("/things/found")
	group.PATCH("/things/found")
	group.GET("/thins/found")

	//二手市场
	group.POST("/market/used")
	group.PATCH("/market/used")
	group.GET("/market/used")

	//匿名树洞
	group.POST("/anonymous")
	group.PATCH("/anonymous")
	group.GET("/anonymous")
	group.DELETE("/anonymous")

	//资源下载
	group.POST("/resources")
	group.GET("/resources")
	group.DELETE("/resources")

	//博文
	group.POST("/blog")
	group.PATCH("/blog")
	group.GET("/blog")
	group.DELETE("/blog")

	//话题讨论
	group.POST("/topic/discussion")
	group.PATCH("/topic/discussion")
	group.GET("/topic/discussion")
	group.DELETE("/topic/discussion")

	//校园咨询
	group.POST("/advisory")
	group.PATCH("/advisory")
}
