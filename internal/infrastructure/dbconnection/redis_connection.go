package dbconnection

import (
	"fmt"
	"points/internal/domain/port"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST" validate:"required"`
	Port     string `mapstructure:"REDIS_PORT" validate:"required"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type RedisConnection struct {
	config   port.Config
	validate *validator.Validate
}

func NewRedisConnection(config port.Config) *RedisConnection {
	return &RedisConnection{config: config, validate: validator.New()}
}

func (d *RedisConnection) InitRedisDatabase() (*redis.Client, error) {
	var cfg RedisConfig
	v := d.config.Sub("redis")
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal redis config: %w", err)
	}

	err := d.validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to validate redis config: %w", err)
	}

	pass := ""
	if cfg.Password != "" {
		pass = ":" + cfg.Password + "@"
	}

	dsn := fmt.Sprintf("redis://%s%s:%s/%d", pass, cfg.Host, cfg.Port, cfg.DB)
	options, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis DSN: %w", err)
	}

	return redis.NewClient(options), nil
}
