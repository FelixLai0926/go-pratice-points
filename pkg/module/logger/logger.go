package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(env string) error {
	var cfg zap.Config
	if strings.ToLower(env) == "production" {
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
	if logger == nil {
		return
	}

	if err := logger.Sync(); err != nil {
		zap.L().Error("Failed to sync logger", zap.Error(err))
	}
}
