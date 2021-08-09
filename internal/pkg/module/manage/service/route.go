package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/manage/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/api/admin/login", handle.AdminLoginHandle)
	app := engine.Group("/api/manage", middleware.SetUserMeta())

	app.GET("/admin/logout", handle.AdminLogoutHandle)
	app.GET("/admin/info/cache", handle.CacheInfoHandle)

	app.GET("/admin/user", handle.UserHandle)

	app.GET("/admin/login/log/admin", handle.AdminLoginLogHandle)
	app.GET("/admin/login/log/user", handle.UserLoginLogHandle)

	app.POST("/admin/shop/commodity", handle.AddCommodityHandle)
	app.GET("/admin/shop/commodity", handle.CommodityHandle)

	app.POST("/admin/coin/giveaway", handle.SystemGiveawayUserCoinHandle)
	app.GET("/admin/coin", handle.UserCoinHandle)

	app.GET("/admin/sms/log", handle.SmsLogHandle)
}
