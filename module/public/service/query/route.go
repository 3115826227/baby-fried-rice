package query

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/public/service/query/handle"
)

func Register(engine *gin.RouterGroup) {

	engine.GET("/train_meta", handle.TrainMetaGet)
}
