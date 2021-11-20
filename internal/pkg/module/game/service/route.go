package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/game/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/api/game", middleware.SetUserMeta())
	app.GET("/record", handle.GameRecordQueryHandle)
	app.GET("/record/detail", handle.GameRecordDetailQueryHandle)

	app.GET("/china_chess/status_data", handle.ChinaChessGameStatusDataQueryHandle)
}
