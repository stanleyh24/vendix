package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

func InitLogger() {
	config := zap.NewProductionConfig()

	// Set log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}

	sugar = log.Sugar()
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func Debug(msg string, fields ...interface{}) {
	sugar.Debugw(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	sugar.Infow(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	sugar.Warnw(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	sugar.Errorw(msg, fields...)
}

func Fatal(msg string, fields ...interface{}) {
	sugar.Fatalw(msg, fields...)
}

func GetLogger() *zap.Logger {
	return log
}

func GetSugaredLogger() *zap.SugaredLogger {
	return sugar
}
