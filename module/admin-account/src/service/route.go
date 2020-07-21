package service

import (
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/service/handle"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(engine *gin.Engine) {

	//登陆
	engine.POST("/api/admin/login", handle.AdminLoginHandle)

	admin := engine.Group("/api/account/admin", middleware.MiddlewareSetUserMeta())

	admin.GET("/logout", handle.AdminLogoutHandle)
	admin.GET("/refresh")

	//信息管理
	admin.PATCH("/password")
	admin.PATCH("/detail")

	//角色管理
	admin.POST("/role")
	admin.PATCH("/role")
	admin.GET("/role")
	admin.DELETE("/role")

	//权限管理
	admin.POST("/permission")
	admin.GET("/permission")
	admin.DELETE("/permission")

	//校园组织
	admin.POST("/school/organize", handle.OrganizeAddHandle)
	admin.PATCH("/school/organize")
	admin.GET("/school/organize", handle.OrganizeGetHandle)
	admin.DELETE("/school/organize")

	//学生信息
	admin.POST("/school/student", handle.StudentAddHandle)
	admin.PATCH("/school/student")
	admin.GET("/school/student", handle.StudentGetHandle)
	admin.DELETE("/school/student")

	//官方群管理
	admin.POST("/school/official/group", handle.OfficialGroupAddHandle)
	admin.GET("/school/official/group", )

	//公告管理
	admin.POST("/school/announcement")
	admin.PATCH("/school/announcement")
	admin.GET("/school/announcement")
	admin.DELETE("/school/announcement")

	//新闻管理
	admin.POST("/school/news")
	admin.PATCH("/school/news")
	admin.GET("/school/news")
	admin.DELETE("/school/news")

	//社团管理
	admin.POST("/school/club")
	admin.PATCH("/school/club")
	admin.GET("/school/club")

	//咨询
	admin.GET("/school/advisory")
}
