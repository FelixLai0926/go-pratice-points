package httputil

import (
	stdErrors "errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"points/errors"
	"points/pkg/models/enum/errcode"
)

func TestFormatError_TableDriven(t *testing.T) {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	tests := []struct {
		name              string
		inputError        error
		message           string
		expectedStatus    int
		expectedStatusStr string
		expectedErrText   string
	}{
		{
			name:              "Non AppError",
			inputError:        stdErrors.New("some standard error"),
			message:           "test non-app error",
			expectedStatus:    http.StatusInternalServerError,
			expectedStatusStr: errcode.ErrInternal.String(),
			expectedErrText:   http.StatusText(http.StatusInternalServerError),
		},
		{
			name: "AppError",
			inputError: &errors.AppError{
				Code: errcode.ErrNotFound,
				Msg:  "not found error",
			},
			message:           "test app error",
			expectedStatus:    http.StatusBadRequest,
			expectedStatusStr: (&errors.AppError{Code: errcode.ErrNotFound}).Code.String(),
			expectedErrText:   http.StatusText(http.StatusBadRequest),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, resp := FormatError(tc.inputError, tc.message, tc.expectedStatus)
			assert.Equal(t, tc.expectedStatus, status)
			assert.Equal(t, tc.expectedStatusStr, resp["status"])
			assert.Equal(t, tc.expectedErrText, resp["error"])
		})
	}
}
