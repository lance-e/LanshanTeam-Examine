package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var ClientLogger *zap.Logger

func init() {
	//配置输出文件
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./userModule.log",
		MaxSize:    10,
		MaxAge:     10,
		MaxBackups: 10,
		LocalTime:  true,
		Compress:   false,
	})
	//配置日志级别
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	//配置日志输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//日志核心记录器
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(os.Stdout, fileWriter),
		level,
	)

	ClientLogger = zap.New(core)
}
