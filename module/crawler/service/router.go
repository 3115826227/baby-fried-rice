package service

import "github.com/gin-gonic/gin"

func Route(engine *gin.Engine) {

	engine.GET("/api/trigger/station", StationTrigger)
	engine.GET("/api/trigger/train/meta", TrainMetaTrigger)
	engine.GET("/api/trigger/train/seat", TrainSeatTrigger)

	engine.GET("/api/find/train/meta", FindTrainMeta)
}
