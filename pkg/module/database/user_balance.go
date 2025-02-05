package database

import (
	"fmt"
	"points/pkg/model"

	"gorm.io/gorm"
)

func GetUserBalance(db *gorm.DB, userID int) (float64, error) {
	var userBalance float64
	err := db.Model(&model.Balance{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&userBalance).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get user balance: %w", err)
	}

	return userBalance, nil
}
