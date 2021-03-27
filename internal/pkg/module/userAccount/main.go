package userAccount

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Main() {
	engine := gin.Default()

	engine.Use(middleware.Cors())
	//service.Register(engine)

	conf := config.GetConfig()
	engine.Run(fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.Port))
}
