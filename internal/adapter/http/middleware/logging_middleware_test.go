package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
			logChecks:  []string{`"method":"GET"`, `"path":"/test"`, `"status":200`, `"client_ip":`},
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
			var logBuffer bytes.Buffer
			writer := zapcore.AddSync(&logBuffer)
			encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			core := zapcore.NewCore(encoder, writer, zap.DebugLevel)
			logger := zap.New(core)

			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.Use(LoggerMiddleware(logger))
			router.GET(tt.route, tt.handler)

			req, _ := http.NewRequest(http.MethodGet, tt.route, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())

			logOutput := logBuffer.String()
			assert.NotEmpty(t, logOutput, "Not empty log output")

			for _, check := range tt.logChecks {
				assert.Contains(t, logOutput, check, "Log output should contain: "+check)
			}

			var logEntries []map[string]interface{}
			lines := strings.Split(logOutput, "\n")
			for _, line := range lines {
				if len(line) == 0 {
					continue
				}
				var entry map[string]interface{}
				if err := json.Unmarshal([]byte(line), &entry); err == nil {
					logEntries = append(logEntries, entry)
				}
			}
			found := false
			for _, entry := range logEntries {
				if m, ok := entry["method"].(string); ok && m == "GET" {
					found = true
					break
				}
			}
			assert.True(t, found, "At least one log entry should contain method=GET")
		})
	}
}
