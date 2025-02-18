package trade

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"points/api/trade/mock"
	"points/errors"
	"points/pkg/middleware"
	"points/pkg/models/enum/errcode"

	"github.com/gin-gonic/gin"
)

func TestTransferHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                     string
		requestBody              map[string]interface{}
		validateErr              error
		transferErr              error
		expectedHTTPStatus       int
		expectedResponseContains string
	}{
		{
			name: "ValidationError",
			requestBody: map[string]interface{}{
				"from":         1,
				"to":           2,
				"nonce":        12345,
				"amount":       100.0,
				"auto_confirm": true,
			},
			validateErr: &errors.AppError{
				Code: errcode.ErrNotFound,
				Msg:  "source account not found",
			},
			transferErr:              nil,
			expectedHTTPStatus:       http.StatusBadRequest,
			expectedResponseContains: errcode.ErrNotFound.String(),
		},
		{
			name: "TransferServiceError",
			requestBody: map[string]interface{}{
				"from":         1,
				"to":           2,
				"nonce":        12345,
				"amount":       100.0,
				"auto_confirm": true,
			},
			validateErr: nil,
			transferErr: &errors.AppError{
				Code: errcode.ErrInternal,
				Msg:  "dummy transfer error",
			},
			expectedHTTPStatus:       http.StatusInternalServerError,
			expectedResponseContains: errcode.ErrInternal.String(),
		},
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"from":         1,
				"to":           2,
				"nonce":        12345,
				"amount":       100.0,
				"auto_confirm": true,
			},
			validateErr:              nil,
			transferErr:              nil,
			expectedHTTPStatus:       http.StatusOK,
			expectedResponseContains: errcode.ErrOK.String(),
		},
		{
			name: "Success(destination account not exists)",
			requestBody: map[string]interface{}{
				"from":         1,
				"to":           999,
				"nonce":        123456,
				"amount":       100.0,
				"auto_confirm": true,
			},
			validateErr:              nil,
			transferErr:              nil,
			expectedHTTPStatus:       http.StatusOK,
			expectedResponseContains: errcode.ErrOK.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.ErrorHandlerMiddleware())
			dummyService := &mock.DummyTradeService{
				TransferErr: tc.transferErr,
			}
			dummyValidator := &mock.DummyTransValidator{
				ValidateErr: tc.validateErr,
			}
			handler := NewTransferHandler(nil, dummyService, dummyValidator)
			router.POST("/transfer", handler.Transfer)

			jsonData, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/transfer", bytes.NewBuffer(jsonData))
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
