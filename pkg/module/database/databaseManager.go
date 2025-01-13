package database

import "gorm.io/gorm"

type DatabaseManager struct {
	DB *gorm.DB
}
