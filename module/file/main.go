package main

import (
	_ "github.com/3115826227/baby-fried-rice/module/file/src/config"
	"github.com/3115826227/baby-fried-rice/module/file/src/middleware"
	_ "github.com/3115826227/baby-fried-rice/module/file/src/model"
	"github.com/3115826227/baby-fried-rice/module/file/src/service"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	engine.Use(middleware.Cors())
	service.RegisterRoute(engine)

	engine.Run(":8051")
}
