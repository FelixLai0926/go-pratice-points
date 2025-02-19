package middleware

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"points/errors"
	"points/pkg/models/enum/errcode"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestErrorHandlerMiddleware_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		route            string
		handler          gin.HandlerFunc
		expectedHTTPCode int
		expectedResponse errorResponse
	}{
		{
			name:  "AppError",
			route: "/app-error",
			handler: func(c *gin.Context) {
				appErr := &errors.AppError{
					Code: errcode.ErrInternal,
				}
				c.Error(appErr)
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: errorResponse{
				Status:  errcode.ErrInternal.String(),
				Message: http.StatusText(http.StatusInternalServerError),
			},
		},
		{
			name:  "GenericError",
			route: "/generic-error",
			handler: func(c *gin.Context) {
				c.Error(stdErrors.New("some generic error"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: errorResponse{
				Status:  errcode.ErrInternal.String(),
				Message: http.StatusText(http.StatusInternalServerError),
			},
		},
	}

	router := gin.New()
	router.Use(ErrorHandlerMiddleware())

	for _, tc := range tests {
		router.GET(tc.route, tc.handler)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.route, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedHTTPCode, w.Code)

			var resp errorResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedResponse.Status, resp.Status)
			assert.Equal(t, tc.expectedResponse.Message, resp.Message)
		})
	}
}
