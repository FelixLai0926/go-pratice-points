package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUserBalance(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer sqlDB.Close()

	type args struct {
		db     *gorm.DB
		userID int
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Valid user balance",
			args: args{
				db:     db,
				userID: 0,
			},
			want:    100.0,
			wantErr: false,
		}, {
			name: "user not found",
			args: args{
				db:     db,
				userID: -1,
			},
			want:    0.0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserBalance(tt.args.db, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
