// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package orm

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameAccount = "account"

// Account mapped from table <account>
type Account struct {
	UserID           int64           `gorm:"column:user_id;primaryKey" json:"user_id"`
	AvailableBalance decimal.Decimal `gorm:"column:available_balance;not null" json:"available_balance"`
	ReservedBalance  decimal.Decimal `gorm:"column:reserved_balance;not null" json:"reserved_balance"`
	UpdatedAt        time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Account's table name
func (*Account) TableName() string {
	return TableNameAccount
}
