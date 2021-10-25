package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/manage/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	handle.InitBackend()

	engine.POST("/api/admin/login", handle.AdminLoginHandle)
	app := engine.Group("/api/manage", middleware.SetUserMeta())

	app.GET("/admin/logout", handle.AdminLogoutHandle)
	app.GET("/admin/info/cache", handle.CacheInfoHandle)

	app.GET("/admin/user", handle.UserHandle)
	app.POST("/admin/user", handle.AddUserHandle)

	app.GET("/admin/login/log/admin", handle.AdminLoginLogHandle)
	app.GET("/admin/login/log/user", handle.UserLoginLogHandle)

	app.POST("/admin/shop/commodity", handle.AddCommodityHandle)
	app.PATCH("/admin/shop/commodity", handle.UpdateCommodityHandle)
	app.GET("/admin/shop/commodity", handle.CommodityHandle)
	app.DELETE("/admin/shop/commodity", handle.DeleteCommodityHandle)

	app.GET("/admin/shop/order", handle.OrderHandle)

	app.GET("/admin/space", handle.SpaceHandle)
	app.PATCH("/admin/space/audit", handle.UpdateSpaceAuditHandle)

	app.POST("/admin/coin/giveaway", handle.SystemGiveawayUserCoinHandle)
	app.GET("/admin/coin", handle.UserCoinHandle)

	app.GET("/admin/sms/log", handle.SmsLogHandle)

	app.POST("/admin/iter/version", handle.AddIterativeVersionHandle)
	app.PATCH("/admin/iter/version", handle.UpdateIterativeVersionHandle)
	app.GET("/admin/iter/version", handle.QueryIterativeVersionHandle)
}
