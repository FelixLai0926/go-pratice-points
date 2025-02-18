// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package orm

import (
	"time"
)

const TableNameTransaction = "transaction"

// Transaction mapped from table <transaction>
type Transaction struct {
	TransactionID string    `gorm:"column:transaction_id;not null;default:gen_random_uuid()" json:"transaction_id"`
	Nonce         int64     `gorm:"column:nonce;primaryKey" json:"nonce"`
	FromAccountID int32     `gorm:"column:from_account_id;primaryKey" json:"from_account_id"`
	ToAccountID   int32     `gorm:"column:to_account_id;not null" json:"to_account_id"`
	Amount        float64   `gorm:"column:amount;not null" json:"amount"`
	Status        int32     `gorm:"column:status;not null" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Transaction's table name
func (*Transaction) TableName() string {
	return TableNameTransaction
}
