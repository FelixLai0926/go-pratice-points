package dbconnection

import (
	"errors"
	"testing"

	"points/test/mock"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestInitRedisDatabase_UnmarshalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("redis").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).Return(errors.New("unmarshal error"))

	conn := NewRedisConnection(mockCfg)
	client, err := conn.InitRedisDatabase()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal redis config")
}

func TestInitRedisDatabase_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidCfg := RedisConfig{
		Host:     "",
		Port:     "6379",
		Password: "secret",
		DB:       0,
	}
	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("redis").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*RedisConfig)) = invalidCfg
		return nil
	})

	conn := NewRedisConnection(mockCfg)
	client, err := conn.InitRedisDatabase()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to validate redis config")
}

func TestInitRedisDatabase_DSNParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	invalidCfg := RedisConfig{
		Host:     "localhost",
		Port:     "invalidport",
		Password: "secret",
		DB:       0,
	}
	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("redis").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*RedisConfig)) = invalidCfg
		return nil
	})

	conn := NewRedisConnection(mockCfg)
	client, err := conn.InitRedisDatabase()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse redis DSN")
}

func TestInitRedisDatabase_Success_WithPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validCfg := RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "secret",
		DB:       0,
	}
	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("redis").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*RedisConfig)) = validCfg
		return nil
	})

	conn := NewRedisConnection(mockCfg)
	client, err := conn.InitRedisDatabase()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, &redis.Client{}, client)
}

func TestInitRedisDatabase_Success_NoPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validCfg := RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}
	mockCfg := mock.NewMockConfig(ctrl)
	dummySettingsManager := mock.NewMockSettingsManager(ctrl)
	mockCfg.EXPECT().Sub("redis").Return(dummySettingsManager).Times(1)
	dummySettingsManager.EXPECT().Unmarshal(gomock.Any()).DoAndReturn(func(out interface{}) error {
		*(out.(*RedisConfig)) = validCfg
		return nil
	})

	conn := NewRedisConnection(mockCfg)
	client, err := conn.InitRedisDatabase()
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, &redis.Client{}, client)
}
