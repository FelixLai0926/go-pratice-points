package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		route      string
		handler    gin.HandlerFunc
		wantStatus int
		wantBody   string
		logChecks  []string
	}{
		{
			name:  "Valid request",
			route: "/test",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"success"}`,
			logChecks:  []string{`"method":"GET"`, `"path":"/test"`, `"status":200`},
		},
		{
			name:  "Request with error",
			route: "/error",
			handler: func(c *gin.Context) {
				c.Error(http.ErrNotSupported)
				c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
			},
			wantStatus: http.StatusNotImplemented,
			wantBody:   `{"error":"not implemented"}`,
			logChecks:  []string{`"method":"GET"`, `"path":"/error"`, `"status":501`, `"error":"feature not supported"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 logger
			var logBuffer bytes.Buffer
			writer := zapcore.AddSync(&logBuffer)
			encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			core := zapcore.NewCore(encoder, writer, zap.DebugLevel)
			logger := zap.New(core)
			zap.ReplaceGlobals(logger)

			// 設置 Gin 路由
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(LoggerMiddleware())

			router.GET(tt.route, tt.handler)

			// 發送測試請求
			req, _ := http.NewRequest(http.MethodGet, tt.route, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// 確認響應
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())

			// 驗證日誌輸出
			logOutput := logBuffer.String()
			for _, check := range tt.logChecks {
				assert.Contains(t, logOutput, check, "Log should contain: "+check)
			}
		})
	}
}
