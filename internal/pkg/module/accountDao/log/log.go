package log

import "baby-fried-rice/internal/pkg/kit/log"

var (
	Logger log.Logging
)

func InitLog(serviceName, logLevel, logPath string) (err error) {
	Logger, err = log.NewLoggerClient(serviceName, logLevel, logPath)
	return
}
