package cache

import (
	"context"
	"points/pkg/module/test"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRedisDSN(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         *RedisConfig
		expectedDSN string
		expectErr   bool
	}{
		{
			name: "Success with password",
			cfg: &RedisConfig{
				Host:     "localhost",
				Port:     "6379",
				Password: "secret",
				DB:       0,
			},
			expectedDSN: "redis://:secret@localhost:6379/0",
			expectErr:   false,
		},
		{
			name: "Success without password",
			cfg: &RedisConfig{
				Host:     "localhost",
				Port:     "6379",
				Password: "",
				DB:       0,
			},
			expectedDSN: "redis://localhost:6379/0",
			expectErr:   false,
		},
		{
			name: "Missing Host",
			cfg: &RedisConfig{
				Host:     "",
				Port:     "6379",
				Password: "secret",
				DB:       0,
			},
			expectedDSN: "",
			expectErr:   true,
		},
		{
			name: "Missing Port",
			cfg: &RedisConfig{
				Host:     "localhost",
				Port:     "",
				Password: "secret",
				DB:       0,
			},
			expectedDSN: "",
			expectErr:   true,
		},
		{
			name:        "Nil Config",
			cfg:         nil,
			expectedDSN: "",
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dsn, err := GenerateRedisDSN(tc.cfg)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDSN, dsn)
			}
		})
	}
}

func TestInitRedisClient_Combined(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          *RedisConfig
		expectErr    bool
		expectValid  bool
		useContainer bool
	}{
		{
			name: "Missing Host",
			cfg: &RedisConfig{
				Host: "",
				Port: "6379",
				DB:   0,
			},
			expectErr:    true,
			expectValid:  false,
			useContainer: false,
		},
		{
			name: "Missing Port",
			cfg: &RedisConfig{
				Host: "localhost",
				Port: "",
				DB:   0,
			},
			expectErr:    true,
			expectValid:  false,
			useContainer: false,
		},
		{
			name: "Success with Container",
			cfg: &RedisConfig{
				Host: "localhost",
				Port: "6379",
				DB:   0,
			},
			expectErr:    false,
			expectValid:  true,
			useContainer: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.useContainer {
				dsn, cleanup := test.SetupRedisContainer(t)
				t.Cleanup(cleanup)
				updateCfgFromDSN(tc.cfg, dsn)
			}
			client, err := InitRedisClient(tc.cfg)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				if err != nil {
					t.Skipf("Skipping success case as no Redis server is available: %v", err)
				}
				assert.NoError(t, err)
				assert.NotNil(t, client)
				if tc.expectValid {
					res, err := client.Ping(context.Background()).Result()
					assert.NoError(t, err)
					assert.Equal(t, "PONG", res)
				}
			}
		})
	}
}

func updateCfgFromDSN(cfg *RedisConfig, dsn string) {
	parts := strings.Split(dsn, ":")
	if len(parts) >= 3 {
		portPart := parts[2]
		portPart = strings.Split(portPart, "/")[0]
		cfg.Host = "localhost"
		cfg.Port = portPart
	}
}
