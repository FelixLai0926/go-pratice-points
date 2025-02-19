package repository

import (
	"context"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"

	"gorm.io/gorm"
)

type TradeRepository interface {
	CreateAccount(ctx context.Context, tx *gorm.DB, userID int64) error
	GetAccount(ctx context.Context, tx *gorm.DB, userID int64) (*orm.Account, error)
	UpdateAccount(ctx context.Context, tx *gorm.DB, account *orm.Account) error
	CreateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error
	CreateOrUpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error
	UpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error
	CreateTransactionEvent(ctx context.Context, tx *gorm.DB, event *orm.TransactionEvent) error
	GetTransaction(ctx context.Context, tx *gorm.DB, nonce, from int64, status *tcc.Status) (*orm.Transaction, error)
}
