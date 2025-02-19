package cache

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" required:"true"`
	Port     string `env:"REDIS_PORT" required:"true"`
	Password string `env:"REDIS_PASSWORD" required:"true"`
	DB       int    `env:"REDIS_DB" required:"true"`
}
