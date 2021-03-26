package service

import (
	"baby-fried-rice/internal/pkg/module/accountDao/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/dao/account")

	app.POST("/root/login", handle.RootLogin)
}
