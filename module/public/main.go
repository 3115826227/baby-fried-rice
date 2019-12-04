package main

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service"
)

func main() {
	engine := gin.Default()

	service.Register(engine)

	engine.Run(":9082")
}
