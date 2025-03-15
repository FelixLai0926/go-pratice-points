package di

import (
	"points/internal/infrastructure/dbconnection"

	"points/internal/domain/port"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var DatabaseModule = fx.Options(
	fx.Provide(
		func(config port.Config) (*gorm.DB, error) {
			pgConn := dbconnection.NewPostgresConnection(config)
			return pgConn.InitPostgresDatabase()
		},
	),
	fx.Provide(
		func(config port.Config) (*redis.Client, error) {
			redisConn := dbconnection.NewRedisConnection(config)
			return redisConn.InitRedisDatabase()
		},
	),
)
