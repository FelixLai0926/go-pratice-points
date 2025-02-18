package repository

import (
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"

	"gorm.io/gorm"
)

type TradeRepository interface {
	CreateAccount(tx *gorm.DB, userID int32) error
	GetAccount(tx *gorm.DB, userID int32) (*orm.Account, error)
	UpdateAccount(tx *gorm.DB, account *orm.Account) error
	CreateTransaction(tx *gorm.DB, trans *orm.Transaction) error
	CreateOrUpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error
	UpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error
	CreateTransactionEvent(tx *gorm.DB, event *orm.TransactionEvent) error
	GetTransaction(tx *gorm.DB, nonce int64, from int32, status *tcc.Status) (*orm.Transaction, error)
}
