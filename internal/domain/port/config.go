package port

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	Sub(key string) SettingsManager
	Unmarshal(out interface{}) error
	SetDefault(req interface{}) error
	SetDefaultInt(key string, value int)
	Copy(dst, src interface{}) error
}
