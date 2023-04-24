package logger

import (
	"os"

	"github.com/shashimalcse/cronuseo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *config.Config) *zap.Logger {

	zap_config := zap.NewProductionEncoderConfig()
	zap_config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zap_config)
	consoleEncoder := zapcore.NewConsoleEncoder(zap_config)
	logFile, _ := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}
