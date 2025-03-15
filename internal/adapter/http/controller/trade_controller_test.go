package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"points/internal/adapter/http/middleware"
	"points/internal/domain/command"
	"points/internal/domain/valueobject"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
	"points/test/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupRouter(route string, method string, handler gin.HandlerFunc) (*gin.Engine, *bytes.Buffer) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	var logBuffer bytes.Buffer
	writer := zapcore.AddSync(&logBuffer)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writer, zap.DebugLevel)
	logger := zap.New(core)

	router.Use(middleware.ErrorHandlerMiddleware(logger))
	router.Handle(method, route, handler)
	return router, &logBuffer
}

func newTestTradeController(ctrl *gomock.Controller) *TradeController {
	mockTradeUsecase := mock.NewMockTradeUsecase(ctrl)
	mockConfig := mock.NewMockConfig(ctrl)
	mockConfig.EXPECT().SetDefault(gomock.Any()).Return(nil).AnyTimes()
	return NewTradeController(mockTradeUsecase, mockConfig)
}

func newExpectedTransferCommand(autoConfirm bool) *command.TransferCommand {
	return &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 12345,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
		AutoConfirm: autoConfirm,
	}
}

func TestTransferHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name                string
		requestBody         string
		expectedAutoConfirm bool
		tradeUsecaseErr     error
		expectedHTTPStatus  int
		expectedResponseStr string
	}{
		{
			name: "Success auto_confirm omitted (defaults to true)",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345,
				"amount": 100.0
			}`,
			expectedAutoConfirm: true,
			tradeUsecaseErr:     nil,
			expectedHTTPStatus:  http.StatusOK,
			expectedResponseStr: errcode.ErrOK.String(),
		},
		{
			name: "Success auto_confirm explicitly true",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345,
				"amount": 100.0,
				"auto_confirm": true
			}`,
			expectedAutoConfirm: true,
			tradeUsecaseErr:     nil,
			expectedHTTPStatus:  http.StatusOK,
			expectedResponseStr: errcode.ErrOK.String(),
		},
		{
			name: "Success auto_confirm explicitly false",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345,
				"amount": 100.0,
				"auto_confirm": false
			}`,
			expectedAutoConfirm: false,
			tradeUsecaseErr:     nil,
			expectedHTTPStatus:  http.StatusOK,
			expectedResponseStr: errcode.ErrOK.String(),
		},
		{
			name:                "Validation Error",
			requestBody:         `{"invalid": "data"}`,
			expectedAutoConfirm: false,
			tradeUsecaseErr:     nil,
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedResponseStr: errcode.ErrInvalidRequest.String(),
		},
		{
			name: "Transfer Service Error",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345,
				"amount": 100.0,
				"auto_confirm": true
			}`,
			expectedAutoConfirm: true,
			tradeUsecaseErr:     apperror.Wrap(errcode.ErrInternal, "dummy transfer error", nil),
			expectedHTTPStatus:  http.StatusInternalServerError,
			expectedResponseStr: errcode.ErrInternal.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tradeController := newTestTradeController(ctrl)
			router, _ := setupRouter("/transfer", http.MethodPost, tradeController.Transfer)

			mockTradeUsecase := tradeController.TradeUsecase.(*mock.MockTradeUsecase)
			mockConfig := tradeController.config.(*mock.MockConfig)
			expectedCmd := newExpectedTransferCommand(tc.expectedAutoConfirm)
			mockConfig.EXPECT().
				Copy(gomock.Any(), gomock.Any()).
				DoAndReturn(func(dest interface{}, src interface{}) error {
					if d, ok := dest.(*command.TransferCommand); ok {
						*d = *expectedCmd
					}
					return nil
				}).AnyTimes()
			mockTradeUsecase.EXPECT().
				Transfer(gomock.Any(), newExpectedTransferCommand(tc.expectedAutoConfirm)).
				Return(tc.tradeUsecaseErr).AnyTimes()

			req, err := http.NewRequest("POST", "/transfer", strings.NewReader(tc.requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedHTTPStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedResponseStr)
		})
	}
}

func TestConfirmHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name                string
		requestBody         string
		confirmErr          error
		expectedHTTPStatus  int
		expectedResponseStr string
	}{
		{
			name: "Success Confirm",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345
			}`,
			confirmErr:          nil,
			expectedHTTPStatus:  http.StatusOK,
			expectedResponseStr: errcode.ErrOK.String(),
		},
		{
			name:                "Validation Error Confirm",
			requestBody:         `{"invalid": "data"}`,
			confirmErr:          nil,
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedResponseStr: errcode.ErrInvalidRequest.String(),
		},
		{
			name: "Confirm Service Error",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345
			}`,
			confirmErr:          apperror.Wrap(errcode.ErrInternal, "dummy confirm error", nil),
			expectedHTTPStatus:  http.StatusInternalServerError,
			expectedResponseStr: errcode.ErrInternal.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tradeController := newTestTradeController(ctrl)
			router, _ := setupRouter("/confirm", http.MethodPost, tradeController.Confirm)

			mockTradeUsecase := tradeController.TradeUsecase.(*mock.MockTradeUsecase)
			mockCopier := tradeController.config.(*mock.MockConfig)

			expectedCmd := &command.ConfirmCommand{
				BaseCommand: command.BaseCommand{
					From:  1,
					To:    2,
					Nonce: 12345,
				},
			}
			mockCopier.EXPECT().
				Copy(gomock.Any(), gomock.Any()).
				DoAndReturn(func(dest interface{}, src interface{}) error {
					if d, ok := dest.(*command.ConfirmCommand); ok {
						*d = *expectedCmd
					}
					return nil
				}).AnyTimes()
			mockTradeUsecase.EXPECT().
				ManualConfirm(gomock.Any(), expectedCmd).
				Return(tc.confirmErr).AnyTimes()

			req, err := http.NewRequest("POST", "/confirm", strings.NewReader(tc.requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedHTTPStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedResponseStr)
		})
	}
}

func TestCancelHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name                string
		requestBody         string
		cancelErr           error
		expectedHTTPStatus  int
		expectedResponseStr string
	}{
		{
			name: "Success Cancel",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345
			}`,
			cancelErr:           nil,
			expectedHTTPStatus:  http.StatusOK,
			expectedResponseStr: errcode.ErrOK.String(),
		},
		{
			name:                "Validation Error Cancel",
			requestBody:         `{"invalid": "data"}`,
			cancelErr:           nil,
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedResponseStr: errcode.ErrInvalidRequest.String(),
		},
		{
			name: "Cancel Service Error",
			requestBody: `{
				"from": 1,
				"to": 2,
				"nonce": 12345
			}`,
			cancelErr:           apperror.Wrap(errcode.ErrInternal, "dummy cancel error", nil),
			expectedHTTPStatus:  http.StatusInternalServerError,
			expectedResponseStr: errcode.ErrInternal.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tradeController := newTestTradeController(ctrl)
			router, _ := setupRouter("/cancel", http.MethodPost, tradeController.Cancel)

			mockTradeUsecase := tradeController.TradeUsecase.(*mock.MockTradeUsecase)
			mockConfig := tradeController.config.(*mock.MockConfig)

			expectedCmd := &command.CancelCommand{
				BaseCommand: command.BaseCommand{
					From:  1,
					To:    2,
					Nonce: 12345,
				},
			}
			mockConfig.EXPECT().
				Copy(gomock.Any(), gomock.Any()).
				DoAndReturn(func(dest interface{}, src interface{}) error {
					if d, ok := dest.(*command.CancelCommand); ok {
						*d = *expectedCmd
					}
					return nil
				}).AnyTimes()
			mockTradeUsecase.EXPECT().
				Cancel(gomock.Any(), expectedCmd).
				Return(tc.cancelErr).AnyTimes()

			req, err := http.NewRequest("POST", "/cancel", strings.NewReader(tc.requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedHTTPStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedResponseStr)
		})
	}
}
