package main

import (
	"github.com/3115826227/baby-fried-rice/module/crawler/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
)

/*
	开启车次信息消费端
*/
func TrainMetaConsumerOpen() {
	for i := 0; i < 2; i++ {
		go service.ZhixingTrainConsumer()
		//go service.QunarTrainConsumer()
		go service.MeituanTrainConsumer()
	}
	//go service.TrainTaskConsumer()
	for i := 0; i < 4; i++ {
		go service.TrainRelInsertConsumer()
	}
	go service.TrainMetaInsertConsumer()
}

/*
	开启列车坐席消费端
*/
func TrainSeatConsumerOpen() {
	for i := 0; i < 4; i++ {
		go service.SeatInsertConsumer()
	}
	for i := 0; i < 2; i++ {
		//go service.TongChengYiLongConsumer()
		go service.ZhiXingConsumer()
		//go service.QunarConsumer()
		go service.MeituanConsumer()
		go service.JindongConsumer()
	}
	go service.SeatConsumer()
}

func init() {
	logPath := os.Getenv("ACCESS_LOG_PATH")
	if logPath != "" {
		logDir := path.Dir(logPath)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			log.Fatal("ERROR 日志目录 ", logDir, " 不存在")
		}

		// 打印到文件，自动分裂
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    64, // megabytes
			MaxBackups: 10,
			MaxAge:     28, // days
		})

		gin.DefaultWriter = w
	}
	TrainMetaConsumerOpen()
	TrainSeatConsumerOpen()
}

func main() {

	engine := gin.Default()

	service.Route(engine)

	engine.Run(":9083")
}
