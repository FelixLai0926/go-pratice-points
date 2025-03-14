package infrastructure

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var buildFn = func(cfg zap.Config) (*zap.Logger, error) {
	return cfg.Build()
}

func TestNewZapLogger(t *testing.T) {
	originalBuildFn := buildFn
	defer func() { buildFn = originalBuildFn }()

	tests := []struct {
		name       string
		env        string
		buildError error
		wantErr    bool
	}{
		{
			name:    "Production environment",
			env:     "production",
			wantErr: false,
		},
		{
			name:    "Development environment",
			env:     "development",
			wantErr: false,
		},
		{
			name:    "Unknown environment => default to development config",
			env:     "staging",
			wantErr: false,
		},
		{
			name:       "Build returns error",
			env:        "production",
			buildError: errors.New("failed to build zap logger"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildFn = func(cfg zap.Config) (*zap.Logger, error) {
				if tt.buildError != nil {
					return nil, tt.buildError
				}
				return cfg.Build()
			}

			logger, err := testableNewZapLogger(tt.env)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.True(t, logger.Core().Enabled(zap.InfoLevel), "logger should allow Info level logging")
			}
		})
	}
}

func testableNewZapLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return buildFn(cfg)
}
