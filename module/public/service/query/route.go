package query

import (
	"github.com/3115826227/baby-fried-rice/module/public/service/query/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.RouterGroup) {

	railway := engine.Group("/railway")
	railway.GET("/city", handle.CityGet)
	railway.GET("/train_meta", handle.TrainMetaGet)
	railway.POST("/train_seat/price", handle.TrainSeatPriceGet)
}
