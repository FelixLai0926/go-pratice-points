package router

import (
	"points/api/test"

	"github.com/gin-gonic/gin"
)

func RegisterTestRoutes(server *gin.Engine) {
	apiTest := server.Group("/test")
	{
		apiTest.GET("/ping", test.Ping)
	}
}
