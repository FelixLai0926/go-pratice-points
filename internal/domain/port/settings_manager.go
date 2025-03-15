package port

type SettingsManager interface {
	SetDefault(key string, value interface{})
	GetString(key string) string
	GetInt(key string) int
	Sub(key string) SettingsManager
	Unmarshal(out interface{}) error
}
