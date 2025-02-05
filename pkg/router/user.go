package router

import (
	"points/api/users"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(server *gin.Engine) {
	user := server.Group("/user")
	{
		user.POST("/transfer", users.Transfer)
	}
}
