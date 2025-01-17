package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		zap.L().Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()

		duration := time.Since(startTime)
		zap.L().Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)

		for _, err := range c.Errors {
			zap.L().Error("Request error", zap.String("error", err.Error()))
		}
	}
}
