package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func (*TestController) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
