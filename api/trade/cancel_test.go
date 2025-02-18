package trade

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"points/api/trade/mock"
	"points/errors"
	"points/pkg/middleware"
	"points/pkg/models/enum/errcode"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCancelHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                     string
		requestBody              map[string]interface{}
		validatorErr             error
		cancelErr                error
		expectedHTTPStatus       int
		expectedResponseContains string
	}{
		{
			name: "Validation Error",
			requestBody: map[string]interface{}{
				"from":  1,
				"to":    2,
				"nonce": 1111,
			},
			validatorErr: &errors.AppError{
				Code: errcode.ErrNotFound,
				Msg:  "source account not found",
			},
			cancelErr:                nil,
			expectedHTTPStatus:       http.StatusBadRequest,
			expectedResponseContains: errcode.ErrNotFound.String(),
		},
		{
			name: "Cancel Service Error",
			requestBody: map[string]interface{}{
				"from":  1,
				"to":    2,
				"nonce": 2222,
			},
			validatorErr: nil,
			cancelErr: &errors.AppError{
				Code: errcode.ErrInternal,
				Msg:  "dummy cancel error",
			},
			expectedHTTPStatus:       http.StatusInternalServerError,
			expectedResponseContains: errcode.ErrInternal.String(),
		},
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"from":  1,
				"to":    2,
				"nonce": 3333,
			},
			validatorErr:             nil,
			cancelErr:                nil,
			expectedHTTPStatus:       http.StatusOK,
			expectedResponseContains: errcode.ErrOK.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.ErrorHandlerMiddleware())
			dummySvc := &mock.DummyTradeService{
				CancelErr: tc.cancelErr,
			}
			dummyValidator := &mock.DummyTransValidator{
				ValidateErr: tc.validatorErr,
			}

			handler := NewTransferHandler(nil, dummySvc, dummyValidator)
			router.POST("/cancel", handler.Cancel)

			jsonData, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/cancel", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("failed to create HTTP request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedHTTPStatus {
				t.Errorf("expected HTTP status %d, got %d, body: %s", tc.expectedHTTPStatus, rr.Code, rr.Body.String())
			}

			if !strings.Contains(rr.Body.String(), tc.expectedResponseContains) {
				t.Errorf("expected response to contain %q, got %s", tc.expectedResponseContains, rr.Body.String())
			}
		})
	}
}
