package middleware

import (
	stdErrors "errors"
	"net/http"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		e := c.Errors[len(c.Errors)-1].Err
		var appErr *apperror.AppError
		if stdErrors.As(e, &appErr) {
			logger.Error("Request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("error", appErr.Error()),
			)
			c.JSON(mapErrorCodeToHTTPStatus(appErr.Code), gin.H{
				"status":  appErr.Code.String(),
				"message": http.StatusText(mapErrorCodeToHTTPStatus(appErr.Code)),
			})
		} else {
			logger.Error("Request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("error", e.Error()),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  errcode.ErrInternal.String(),
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		c.Abort()
		return
	}
}
