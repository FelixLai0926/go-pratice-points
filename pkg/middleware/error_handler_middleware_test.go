package middleware

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"net/http/httptest"
	"points/errors"
	"testing"

	"github.com/gin-gonic/gin"
)

func dummyAppErrorHandler(c *gin.Context) {
	appErr := errors.NewAppError(http.StatusBadRequest, stdErrors.New("underlying error"))
	c.Error(appErr)
}

func dummyGenericErrorHandler(c *gin.Context) {
	err := stdErrors.New("generic error")
	c.Error(err)
}

func TestErrorHandlerMiddleware_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandlerMiddleware())
	router.GET("/app-error", dummyAppErrorHandler)

	req, err := http.NewRequest("GET", "/app-error", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var jsonResp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &jsonResp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	expectedMsg := http.StatusText(http.StatusBadRequest)
	if jsonResp["message"] != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, jsonResp["message"])
	}
}

func TestErrorHandlerMiddleware_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandlerMiddleware())
	router.GET("/generic-error", dummyGenericErrorHandler)

	req, err := http.NewRequest("GET", "/generic-error", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	var jsonResp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &jsonResp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	expectedMsg := http.StatusText(http.StatusInternalServerError) // "Internal Server Error"
	if jsonResp["message"] != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, jsonResp["message"])
	}
}
