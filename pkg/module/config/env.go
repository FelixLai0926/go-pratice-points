package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
)

func GetRootPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found, can't locate project root")
		}
		dir = parent
	}
}

func GetEnvPath(environment ...string) (string, error) {
	env := "example"
	if len(environment) > 0 && environment[0] != "" {
		env = environment[0]
	}

	rootPath, err := GetRootPath()
	if err != nil {
		return "", fmt.Errorf("failed to get root path while generating env path: %v", err)
	}

	envFilePath := filepath.Join(rootPath, "configs", fmt.Sprintf(".env.%s", env))

	return envFilePath, nil
}

func InitEnv(envFilePaths ...string) error {
	if err := godotenv.Overload(envFilePaths...); err != nil {
		return fmt.Errorf("failed to load .env file(s): %v, file(s): %v", err, envFilePaths)
	}

	return nil
}

func ParseEnv[TResponse any]() (*TResponse, error) {
	var cfg TResponse

	val := reflect.ValueOf(&cfg).Elem()
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ParseEnv expects a struct type, got: %s", val.Kind())
	}

	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envValue, _ := GetEnvOrDefault(envKey, field.Tag.Get("default"), toString)
		if field.Tag.Get("required") == "true" && envValue == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", envKey)
		}

		if err := setFieldValue(fieldValue, envValue); err != nil {
			return nil, fmt.Errorf("failed to set value for field %s: %v", field.Name, err)
		}
	}

	return &cfg, nil
}

func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		field.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %s", value)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func GetEnvOrDefault[T any](key string, defaultValue T, parser func(string) (T, error)) (T, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	v, err := parser(value)
	if err != nil {
		return defaultValue, err
	}

	return v, nil
}

func toString(value string) (string, error) {
	return value, nil
}
