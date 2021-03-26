package service

import (
	"baby-fried-rice/internal/pkg/module/rootAccount/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/api/root/login", handle.RootLogin)
}
