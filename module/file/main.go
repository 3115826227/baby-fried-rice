package main

import (
	_ "github.com/3115826227/baby-fried-rice/module/file/src/config"
	_ "github.com/3115826227/baby-fried-rice/module/file/src/model"
	"github.com/3115826227/baby-fried-rice/module/file/src/service"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	service.RegisterRoute(engine)

	engine.Run(":8051")
}
