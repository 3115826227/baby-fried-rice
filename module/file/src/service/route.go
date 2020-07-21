package service

import (
	"github.com/3115826227/baby-fried-rice/module/file/src/log"
	"github.com/3115826227/baby-fried-rice/module/file/src/model"
	"github.com/3115826227/baby-fried-rice/module/file/src/model/db"
	"github.com/3115826227/baby-fried-rice/module/file/src/service/handle"
	"github.com/gin-gonic/gin"
	"time"
)

func init() {
	//initOssMeta()
}

func initOssMeta() {

	var now = time.Now()
	var beans = make([]interface{}, 0)
	beans = append(beans, &model.OssMeta{
		CommonIntField: model.CommonIntField{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Domain: "http://q98fgnvaf.bkt.clouddn.com",
		Bucket: "baby-fried-rice",
		Size:   0,
	})
	beans = append(beans, &model.OssMeta{
		CommonIntField: model.CommonIntField{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Domain: "http://qa5t456ve.bkt.clouddn.com",
		Bucket: "babycampus",
		Size:   0,
	})

	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
	}
}

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/api")
	group.POST("/upload", handle.UploadHandle)
	group.GET("/download")
}
