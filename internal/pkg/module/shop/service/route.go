package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/shop/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/api/shop", middleware.SetUserMeta())

	app.GET("/commodity", handle.CommodityQueryHandle)
	app.GET("/commodity/detail", handle.CommodityDetailQueryHandle)

	app.POST("/commodity/order", handle.CommodityOrderAddHandle)
	app.GET("/commodity/order", handle.CommodityOrderQueryHandle)
	app.PATCH("/commodity/order/pay", handle.CommodityOrderPayHandle)
	app.GET("/commodity/order/detail", handle.CommodityOrderDetailQueryHandle)
	app.DELETE("/commodity/order", handle.CommodityOrderDeleteHandle)
}
