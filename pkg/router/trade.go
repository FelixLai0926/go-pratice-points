package router

import (
	"points/api/trade"
	"points/pkg/module/distributedlock"
	"points/repository"
	"points/service/tcc"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterUserRoutes(server *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	balanceRepo := repository.NewTradeRepo()
	lockClient := distributedlock.NewRedisLockClient(redisClient)
	tradeValidator := tcc.NewTCCValidator(db, balanceRepo)
	tradeService := tcc.NewTCCService(db, balanceRepo, tradeValidator, lockClient)
	transferHandler := trade.NewTransferHandler(db, tradeService)

	user := server.Group("/trade")
	{
		user.POST("/transfer", transferHandler.Transfer)
		user.POST("/confirm", transferHandler.Confirm)
		user.POST("/cancel", transferHandler.Cancel)
	}
}
