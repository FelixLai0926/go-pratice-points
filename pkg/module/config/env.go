package config

import (
	"fmt"
	"os"
	"points/configs"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnvData(environment string) (string, error) {
	fileName := fmt.Sprintf("%s.env", strings.ToLower(environment))
	data, err := configs.EnvFS.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read env file %s: %w", fileName, err)
	}
	return string(data), nil
}

func InitEmbeddedEnv(envData string) error {
	envMap, err := godotenv.Parse(strings.NewReader(envData))
	if err != nil {
		return fmt.Errorf("failed to parse env data: %w", err)
	}

	for key, value := range envMap {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set env %s: %w", key, err)
		}
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
			return nil, fmt.Errorf("failed to set value for field %s: %w", field.Name, err)
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
