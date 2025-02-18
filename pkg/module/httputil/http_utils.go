package httputil

import (
	stdErrors "errors"
	"net/http"

	"points/errors"
	"points/pkg/models/enum/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FormatError(err error, message string, statusCode int) (int, gin.H) {
	var appErr *errors.AppError
	if ok := stdErrors.As(err, &appErr); !ok {
		zap.L().Error(message, zap.Error(err))
		return http.StatusInternalServerError, gin.H{
			"status": errcode.ErrInternal.String(),
			"error":  http.StatusText(http.StatusInternalServerError),
		}
	} else {
		zap.L().Error(message, zap.Error(err))
		return statusCode, gin.H{
			"status": appErr.Code.String(),
			"error":  http.StatusText(statusCode),
		}
	}
}
