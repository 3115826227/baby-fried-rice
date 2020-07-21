package main

import (
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
)

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
}

func main() {

	engine := gin.Default()

	service.RegisterRoute(engine)

	engine.Run(":8075")
}
