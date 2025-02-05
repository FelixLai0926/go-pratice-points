package users

import (
	"fmt"
	"net/http"
	"points/pkg/module/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Transfer(context *gin.Context) {
	db := context.MustGet("db").(*gorm.DB)
	var request struct {
		Form *int `json:"form" form:"form" binding:"required"`
	}

	if err := context.ShouldBind(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userBalance, err := database.GetUserBalance(db, *request.Form)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get user balance: %s", err.Error())})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"balance": userBalance,
	})
}
