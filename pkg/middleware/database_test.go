package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDatabaseMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		db         *gorm.DB
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid database connection",
			db:         &gorm.DB{},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"success"}`,
		},
		{
			name:       "Nil database connection",
			db:         nil,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"database connection is not initialized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.Use(DatabaseMiddleware(tt.db))

			router.GET("/test", func(c *gin.Context) {
				db, exists := c.Get("db")
				if !exists || db == nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection is not initialized"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
