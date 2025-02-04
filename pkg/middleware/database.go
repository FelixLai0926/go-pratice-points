package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			c.JSON(500, gin.H{"error": "database connection is not initialized"})
			c.Abort()
			return
		}

		c.Set("db", db)
		c.Next()
	}
}
