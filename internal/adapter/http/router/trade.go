package router

import (
	"points/internal/adapter/http/controller"
	"points/internal/domain/port"
	"points/internal/infrastructure/distributedlock"
	"points/internal/infrastructure/persistence/repository"
	"points/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterUserRoutes(server *gin.Engine, db *gorm.DB, redisClient *redis.Client, config port.Config) {
	unitOfWork := repository.NewGormUnitOfWorkImpl(db, config)
	locker := distributedlock.NewRedisLocker(redisClient)
	tradeUsecase := usecase.NewTradeUsecase(unitOfWork, locker, config)
	tradeController := controller.NewTradeController(tradeUsecase, config)

	user := server.Group("/trade")
	{
		user.POST("/transfer", tradeController.Transfer)
		user.POST("/confirm", tradeController.Confirm)
		user.POST("/cancel", tradeController.Cancel)
	}
}
