package main

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/gateway/config"
	"github.com/3115826227/baby-fried-rice/module/gateway/service"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	service.RegisterRouter(engine)

	fmt.Println(config.Config.Redis)

	engine.Run(":9080")
}
