package infrastructure

import (
	"fmt"
	"points/internal/domain/port"

	"github.com/spf13/viper"
)

type ViperImpl struct {
	viper *viper.Viper
}

var _ port.SettingsManager = (*ViperImpl)(nil)

func NewViperImpl(environment string) (port.SettingsManager, error) {
	v := viper.New()
	v.SetConfigName(environment)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return &ViperImpl{
		viper: v,
	}, nil
}

func (v *ViperImpl) SetDefault(key string, value interface{}) {
	v.viper.SetDefault(key, value)
}

func (v *ViperImpl) SetDefaultInt(key string, value int) {
	v.viper.SetDefault(key, value)
}

func (v *ViperImpl) GetInt(key string) int {
	return v.viper.GetInt(key)
}

func (v *ViperImpl) GetString(key string) string {
	return v.viper.GetString(key)
}

func (v *ViperImpl) Sub(key string) port.SettingsManager {
	sub := v.viper.Sub(key)
	if sub == nil {
		return nil
	}
	return &ViperImpl{
		viper: sub,
	}
}

func (v *ViperImpl) Unmarshal(out interface{}) error {
	return v.viper.Unmarshal(out)
}
