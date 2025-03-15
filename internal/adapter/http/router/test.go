package router

import (
	"points/internal/adapter/http/controller"

	"github.com/gin-gonic/gin"
)

func RegisterTestRoutes(server *gin.Engine) {
	controller := controller.TestController{}
	apiTest := server.Group("/test")
	{
		apiTest.GET("/ping", controller.Ping)
	}
}
