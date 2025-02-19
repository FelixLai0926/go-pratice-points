package middleware

import (
	stdErrors "errors"
	"net/http"
	"points/errors"
	"points/pkg/models/enum/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		e := c.Errors[len(c.Errors)-1].Err
		var appErr *errors.AppError
		if stdErrors.As(e, &appErr) {
			zap.L().Error("Request error",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("error", appErr.Error()),
			)
			c.JSON(appErr.Code.HTTPCode(), gin.H{
				"status":  appErr.Code.String(),
				"message": http.StatusText(appErr.Code.HTTPCode()),
			})
		} else {
			zap.L().Error("Request error",
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
