package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	dsn, err := GenerateRedisDSN(cfg)
	if err != nil {
		return nil, err
	}

	options, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis DSN: %w", err)
	}

	client := redis.NewClient(options)

	return client, nil
}

func GenerateRedisDSN(cfg *RedisConfig) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("RedisConfig is nil")
	}
	if cfg.Host == "" || cfg.Port == "" {
		return "", fmt.Errorf("missing required fields: Host and Port")
	}

	pass := ""
	if cfg.Password != "" {
		pass = ":" + cfg.Password + "@"
	}
	dsn := fmt.Sprintf("redis://%s%s:%s/%d", pass, cfg.Host, cfg.Port, cfg.DB)
	return dsn, nil
}
