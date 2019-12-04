package main

import (
	"github.com/3115826227/baby-fried-rice/module/account/service"
	_ "github.com/3115826227/baby-fried-rice/module/account/service/model"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	service.RegisterRoute(engine)

	engine.Run(":9081")
}
