package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type Logging interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

const (
	LogPathEnvName = "LOG_PATH"
	LogLevel       = "LOG_LEVEL"

	DebugLog = "DEBUG"
	InfoLog  = "INFO"
	WarnLog  = "WARN"
	ErrorLog = "ERROR"
)

var (
	Logger Logging
	levels = map[string]zapcore.Level{
		DebugLog: zap.DebugLevel,
		InfoLog:  zap.InfoLevel,
		WarnLog:  zap.WarnLevel,
		ErrorLog: zap.ErrorLevel,
	}
)

func init() {
	var err error
	Logger, err = NewLoggerClient("global", DebugLog, "")
	if err != nil {
		panic(err)
	}
}

func logLevels() []string {
	return []string{
		DebugLog,
		InfoLog,
		WarnLog,
		ErrorLog,
	}
}

type LoggerClient struct {
	logLevel     string
	serviceName  string
	rootLogger   *zap.Logger
	levelLoggers map[string]*zap.Logger
}

func NewLoggerClient(serviceName, logLevel, logPath string) (client Logging, err error) {
	rootLogger, err := newZapLogger(logLevel, logPath)
	if err != nil {
		return
	}
	levelLoggers := make(map[string]*zap.Logger)
	for _, level := range logLevels() {
		levelLoggers[level] = rootLogger
	}
	client = &LoggerClient{
		logLevel:     logLevel,
		serviceName:  serviceName,
		rootLogger:   rootLogger,
		levelLoggers: levelLoggers,
	}
	return
}

func (lc *LoggerClient) Debug(msg string, fields ...zap.Field) {
	lc.log(DebugLog, msg, fields...)
}

func (lc *LoggerClient) Info(msg string, fields ...zap.Field) {
	lc.log(InfoLog, msg, fields...)
}

func (lc *LoggerClient) Warn(msg string, fields ...zap.Field) {
	lc.log(WarnLog, msg, fields...)
}

func (lc *LoggerClient) Error(msg string, fields ...zap.Field) {
	lc.log(ErrorLog, msg, fields...)
}

func (lc *LoggerClient) log(level string, msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	paths := strings.Split(funcName, "/")
	pack := strings.Split(funcName[strings.LastIndexByte(funcName, '/')+1:], ".")[0]
	paths = append(paths[:len(paths)-1], pack, strings.Split(file[strings.LastIndexByte(file, '/')+1:], ".")[0])
	funcPath := strings.Join(paths, "/")
	caller := funcPath + ".go:" + strconv.Itoa(line)
	now := time.Now().Format("2006/01/02 15:04:05")
	debugStr := Yellow + "[debug] " + Reset
	debugCaller := fmt.Sprintf(" %v%v%v ", Yellow, caller, Reset)
	infoStr := Green + "[info] " + Reset
	infoCaller := fmt.Sprintf(" %v%v%v ", Green, caller, Reset)
	warnStr := Magenta + "[warn] " + Reset
	warnCaller := fmt.Sprintf(" %v%v%v ", Magenta, caller, Reset)
	errStr := Red + "[error] " + Reset
	errCaller := fmt.Sprintf(" %v%v%v ", Red, caller, Reset)
	var message = fmt.Sprintf("%v %v[%v]%v", now, Blue, lc.serviceName, Reset)

	// 日志输出
	switch level {
	case DebugLog:
		lc.levelLoggers[level].Debug(debugStr+message+debugCaller+msg, fields...)
	case InfoLog:
		lc.levelLoggers[level].Info(infoStr+message+infoCaller+msg, fields...)
	case WarnLog:
		lc.levelLoggers[level].Warn(warnStr+message+warnCaller+msg, fields...)
	case ErrorLog:
		lc.levelLoggers[level].Error(errStr+message+errCaller+msg, fields...)
	}
}

func newZapLogger(logLevelStr, logPath string) (zapLog *zap.Logger, err error) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	//encoderConfig := zap.NewProductionEncoderConfig()
	// 选择自定义日志样式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var logLevel zapcore.Level
	var exists bool
	logLevel, exists = levels[logLevelStr]
	if !exists {
		logLevel = zap.DebugLevel
	}

	if logPath != "" {
		logDir := path.Dir(logPath)
		if _, err = os.Stat(logDir); os.IsNotExist(err) {
			log.Fatal("ERROR 日志目录 ", logDir, " 不存在")
			return
		}
		// 打印到文件，自动分裂
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    64, // megabytes
			MaxBackups: 10,
			MaxAge:     28, // days
		})
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			zap.NewAtomicLevelAt(logLevel),
		)
		zapLog = zap.New(core, zap.AddCaller())
	} else {
		// 打印到控制台
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(logLevel)
		cfg.Encoding = "console"
		cfg.EncoderConfig = encoderConfig
		//cfg.DisableStacktrace = true
		zapLog, err = cfg.Build()
		if err != nil {
			log.Fatal("ERROR ", err)
			return
		}
	}
	return
}
