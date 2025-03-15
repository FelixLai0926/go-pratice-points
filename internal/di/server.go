package di

import (
	"context"
	"fmt"
	"points/internal/adapter/http/middleware"
	"points/internal/adapter/http/router"
	"points/internal/domain/port"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var HTTPModule = fx.Options(
	fx.Provide(NewGinServer),
	fx.Invoke(RegisterRoutes),
)

func NewGinServer(config port.Config, logger *zap.Logger) *gin.Engine {
	server := gin.Default()
	server.Use(middleware.LoggerMiddleware(logger))
	server.Use(middleware.ErrorHandlerMiddleware(logger))
	return server
}

func RegisterRoutes(
	server *gin.Engine,
	db *gorm.DB,
	redisClient *redis.Client,
	config port.Config,
) {
	router.RegisterTestRoutes(server)
	router.RegisterUserRoutes(server, db, redisClient, config)
}

func StartServer(lifecycle fx.Lifecycle, server *gin.Engine, config port.Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serverParameter := config.Sub("server")
			addr := fmt.Sprintf("%s:%s", serverParameter.GetString("SERVER_HOST"), serverParameter.GetString("SERVER_PORT"))
			fmt.Println("Starting server...")
			println("test: ", addr)
			go func() {
				if err := server.Run(addr); err != nil {
					panic(fmt.Sprintf("Server failed to start: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return nil
		},
	})
}
