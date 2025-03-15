package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		if c.Request.Body == nil {
			c.Request.Body = http.NoBody
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("read request body error", zap.String("error", err.Error()))
		}

		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("ctx", fmt.Sprintf("%p", c.Request.Context())),
			zap.String("request body", string(bodyBytes)),
		)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Next()

		duration := time.Since(startTime)
		logger.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)

		for _, err := range c.Errors {
			logger.Error("Request error", zap.String("error", err.Error()))
		}
	}
}
