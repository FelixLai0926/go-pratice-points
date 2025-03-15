package di

import (
	"points/internal/infrastructure"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var LoggerModule = fx.Options(
	fx.Provide(func(env string) (*zap.Logger, error) {
		return infrastructure.NewZapLogger(env)
	}),
)
