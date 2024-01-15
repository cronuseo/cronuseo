package logger

import (
	"os"

	"github.com/shashimalcse/cronuseo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(cfg *config.Config) (*zap.Logger, error) {

	zap_config := zap.NewProductionEncoderConfig()
	zap_config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(zap_config)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel)
	if cfg.Log.Enabled {
		fileEncoder := zapcore.NewJSONEncoder(zap_config)
		logDir := "./log"
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return nil, err
		}
		logFile, err := os.OpenFile("./log/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), defaultLogLevel)
		core = zapcore.NewTee(core, fileCore)
	}
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil

}
