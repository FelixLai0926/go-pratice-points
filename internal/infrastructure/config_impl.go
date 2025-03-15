package infrastructure

import (
	"points/internal/domain/port"
)

type ConfigImpl struct {
	seetingManager port.SettingsManager
	defaults       port.Defaults
	copier         port.Copier
}

var _ port.Config = (*ConfigImpl)(nil)

func NewConfigImpl(seetingManager port.SettingsManager, defaults port.Defaults, copier port.Copier) port.Config {
	return &ConfigImpl{
		seetingManager: seetingManager,
		defaults:       defaults,
		copier:         copier,
	}
}

func (c *ConfigImpl) GetString(key string) string {
	return c.seetingManager.GetString(key)
}

func (c *ConfigImpl) GetInt(key string) int {
	return c.seetingManager.GetInt(key)
}

func (c *ConfigImpl) Sub(key string) port.SettingsManager {
	return c.seetingManager.Sub(key)
}

func (c *ConfigImpl) Unmarshal(out interface{}) error {
	return c.seetingManager.Unmarshal(out)
}

func (c *ConfigImpl) SetDefault(req interface{}) error {
	return c.defaults.Set(req)
}

func (c *ConfigImpl) SetDefaultInt(key string, value int) {
	c.seetingManager.SetDefault(key, value)
}

func (c *ConfigImpl) Copy(dst, src interface{}) error {
	return c.copier.Copy(dst, src)
}
