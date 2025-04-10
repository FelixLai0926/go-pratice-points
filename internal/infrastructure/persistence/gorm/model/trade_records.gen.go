// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"github.com/shopspring/decimal"
)

const TableNameTradeRecord = "trade_records"

// TradeRecord mapped from table <trade_records>
type TradeRecord struct {
	TransactionID string          `gorm:"column:transaction_id;not null;default:gen_random_uuid()" json:"transaction_id"`
	Nonce         int64           `gorm:"column:nonce;primaryKey" json:"nonce"`
	FromAccountID int64           `gorm:"column:from_account_id;primaryKey" json:"from_account_id"`
	ToAccountID   int64           `gorm:"column:to_account_id;not null" json:"to_account_id"`
	Amount        decimal.Decimal `gorm:"column:amount;not null" json:"amount"`
	Status        int32           `gorm:"column:status;not null" json:"status"`
	CreatedAt     time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName TradeRecord's table name
func (*TradeRecord) TableName() string {
	return TableNameTradeRecord
}
