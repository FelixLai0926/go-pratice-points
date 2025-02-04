package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(env string) error {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	if logger, err = cfg.Build(); err != nil {
		return fmt.Errorf("InitLogger() error: %v", err)
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func SyncLogger() {
	if logger != nil {
		_ = logger.Sync()
	}
}
