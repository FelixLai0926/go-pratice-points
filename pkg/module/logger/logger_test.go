package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want zapcore.Level
	}{
		{
			name: "Production environment",
			env:  "production",
			want: zapcore.InfoLevel,
		},
		{
			name: "Development environment",
			env:  "development",
			want: zapcore.DebugLevel,
		},
		{
			name: "Unknown environment (defaults to development)",
			env:  "unknown",
			want: zapcore.DebugLevel,
		},
		{
			name: "Uppercase production environment",
			env:  "PRODUCTION",
			want: zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitLogger(tt.env)
			assert.NoError(t, err, "InitLogger() should not return an error")
			assert.NotNil(t, logger, "Logger should not be nil after initialization")

			core := logger.Core()
			enabled := core.Enabled(tt.want)
			assert.True(t, enabled, "Logger should have the expected log level %v", tt.want)
		})
	}
}

func TestSyncLogger(t *testing.T) {
	t.Run("SyncLogger when logger is initialized", func(t *testing.T) {
		err := InitLogger("development")
		assert.NoError(t, err, "InitLogger() should not return an error")

		assert.NotPanics(t, func() {
			SyncLogger()
		}, "SyncLogger should not panic when logger is initialized")
	})

	t.Run("SyncLogger when logger is nil", func(t *testing.T) {
		logger = nil
		assert.NotPanics(t, func() {
			SyncLogger()
		}, "SyncLogger should not panic when logger is nil")
	})
}
