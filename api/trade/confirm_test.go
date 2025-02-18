package trade

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"points/api/trade/mock"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestConfirmHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                     string
		requestBody              map[string]interface{}
		validateErr              error
		manualConfirmErr         error
		expectedHTTPStatus       int
		expectedResponseContains string
	}{
		{
			name: "ValidationError",
			requestBody: map[string]interface{}{
				"from":  999,
				"to":    2,
				"nonce": 123,
			},
			validateErr: &errors.AppError{
				Code: errcode.ErrNotFound,
				Msg:  "source account not found",
			},
			manualConfirmErr:         nil,
			expectedHTTPStatus:       http.StatusBadRequest,
			expectedResponseContains: errcode.ErrNotFound.String(),
		},
		{
			name: "ServiceError",
			requestBody: map[string]interface{}{
				"from":  1,
				"to":    2,
				"nonce": 123,
			},
			validateErr: nil,
			manualConfirmErr: &errors.AppError{
				Code: errcode.ErrInternal,
				Msg:  "dummy manual confirm error",
			},
			expectedHTTPStatus:       http.StatusInternalServerError,
			expectedResponseContains: errcode.ErrInternal.String(),
		},
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"from":  1,
				"to":    2,
				"nonce": 123,
			},
			validateErr:              nil,
			manualConfirmErr:         nil,
			expectedHTTPStatus:       http.StatusOK,
			expectedResponseContains: errcode.ErrOK.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			dummySvc := &mock.DummyTradeService{
				ManualConfirmErr: tc.manualConfirmErr,
			}
			dummyValidator := &mock.DummyTransValidator{
				ValidateErr: tc.validateErr,
			}

			handler := NewTransferHandler(nil, dummySvc, dummyValidator)
			router.POST("/confirm", handler.Confirm)

			jsonData, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/confirm", bytes.NewBuffer(jsonData))
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
