package di

import (
	"points/internal/domain/port"
	"points/internal/infrastructure"

	"go.uber.org/fx"
)

var SettingManagerModule = fx.Options(
	fx.Provide(func(env string) (port.SettingsManager, error) {
		return infrastructure.NewViperImpl(env)
	}),
)

var DefaultsModule = fx.Options(
	fx.Provide(func() port.Defaults {
		return infrastructure.NewDefaultSetterImpl()
	}),
)

var CopierModule = fx.Options(
	fx.Provide(func() port.Copier {
		return infrastructure.NewCopierImpl()
	}),
)

var ConfigModule = fx.Options(
	fx.Provide(func(settingManager port.SettingsManager, defaults port.Defaults, copier port.Copier) port.Config {
		return infrastructure.NewConfigImpl(settingManager, defaults, copier)
	}),
)
