package test

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init() *zap.Logger {

	// cfg, err := config.Load("../../config/local-debug.yml")
	// if err != nil {
	// 	log.Fatal("Error while loading config for test.")
	// 	os.Exit(-1)
	// }
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
