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
			zap.L().Error(appErr.Msg, zap.Error(appErr))
			c.JSON(appErr.HTTPCode, gin.H{
				"status":  appErr.Code.String(),
				"message": http.StatusText(appErr.HTTPCode),
			})
		} else {
			zap.L().Error(http.StatusText(http.StatusInternalServerError), zap.Error(e))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  errcode.ErrInternal.String(),
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}
		c.Abort()
		return
	}
}
