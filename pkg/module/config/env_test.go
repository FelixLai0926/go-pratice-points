package config

import (
	"os"
	"strings"
	"testing"
)

type TestConfig struct {
	Host  string `env:"HOST" default:"localhost" required:"true"`
	Port  int    `env:"PORT" default:"8080"`
	Debug bool   `env:"DEBUG" default:"true"`
}

func TestInitEmbeddedEnv(t *testing.T) {
	os.Unsetenv("TEST_KEY")
	envData := "TEST_KEY=hello\nANOTHER_KEY=world"
	if err := InitEmbeddedEnv(envData); err != nil {
		t.Fatalf("InitEmbeddedEnv error: %v", err)
	}

	if got := os.Getenv("TEST_KEY"); got != "hello" {
		t.Errorf("Expected TEST_KEY to be 'hello', got '%s'", got)
	}
	if got := os.Getenv("ANOTHER_KEY"); got != "world" {
		t.Errorf("Expected ANOTHER_KEY to be 'world', got '%s'", got)
	}

	os.Unsetenv("TEST_KEY")
	os.Unsetenv("ANOTHER_KEY")
}

func TestParseEnvWithDefaults(t *testing.T) {
	os.Unsetenv("HOST")
	os.Unsetenv("PORT")
	os.Unsetenv("DEBUG")

	cfg, err := ParseEnv[TestConfig]()
	if err != nil {
		t.Fatalf("ParseEnv error: %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
	}
	if cfg.Debug != true {
		t.Errorf("Expected Debug to be true, got %v", cfg.Debug)
	}
}

func TestParseEnvWithEnvVariables(t *testing.T) {
	os.Setenv("HOST", "example.com")
	os.Setenv("PORT", "3000")
	os.Setenv("DEBUG", "false")
	defer func() {
		os.Unsetenv("HOST")
		os.Unsetenv("PORT")
		os.Unsetenv("DEBUG")
	}()

	cfg, err := ParseEnv[TestConfig]()
	if err != nil {
		t.Fatalf("ParseEnv error: %v", err)
	}

	if cfg.Host != "example.com" {
		t.Errorf("Expected Host to be 'example.com', got '%s'", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("Expected Port to be 3000, got %d", cfg.Port)
	}
	if cfg.Debug != false {
		t.Errorf("Expected Debug to be false, got %v", cfg.Debug)
	}
}

func TestParseEnvMissingRequired(t *testing.T) {
	type TestRequired struct {
		Value string `env:"MISSING_REQUIRED" required:"true"`
	}
	os.Unsetenv("MISSING_REQUIRED")

	_, err := ParseEnv[TestRequired]()
	if err == nil {
		t.Fatalf("Expected error for missing required environment variable, got nil")
	}
	if !strings.Contains(err.Error(), "missing required environment variable") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_VAR", "set_value")
	defer os.Unsetenv("TEST_VAR")
	value, err := GetEnvOrDefault("TEST_VAR", "default_value", func(s string) (string, error) {
		return s, nil
	})
	if err != nil {
		t.Fatalf("GetEnvOrDefault error: %v", err)
	}
	if value != "set_value" {
		t.Errorf("Expected value to be 'set_value', got '%s'", value)
	}

	os.Unsetenv("TEST_VAR_NOT_SET")
	value, err = GetEnvOrDefault("TEST_VAR_NOT_SET", "default_value", func(s string) (string, error) {
		return s, nil
	})
	if err != nil {
		t.Fatalf("GetEnvOrDefault error: %v", err)
	}
	if value != "default_value" {
		t.Errorf("Expected value to be 'default_value', got '%s'", value)
	}
}
