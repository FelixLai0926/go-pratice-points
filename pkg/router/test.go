package router

import (
	"points/api/test"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	server := gin.Default()

	apiTest := server.Group("/test")
	{
		apiTest.GET("/ping", test.Ping)
	}

	return server
}
