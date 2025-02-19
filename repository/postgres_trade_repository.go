package repository

import (
	"context"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tradeRepo struct{}

func NewTradeRepo() TradeRepository {
	return &tradeRepo{}
}

func (r *tradeRepo) CreateAccount(ctx context.Context, tx *gorm.DB, userID int64) error {
	return tx.WithContext(ctx).Create(&orm.Account{
		UserID:           userID,
		AvailableBalance: decimal.Zero,
		ReservedBalance:  decimal.Zero,
	}).Error
}

func (r *tradeRepo) GetAccount(ctx context.Context, tx *gorm.DB, userID int64) (*orm.Account, error) {
	var account orm.Account
	err := tx.WithContext(ctx).
		Where(&orm.Account{UserID: userID}).
		First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *tradeRepo) UpdateAccount(ctx context.Context, tx *gorm.DB, account *orm.Account) error {
	return tx.WithContext(ctx).Save(account).Error
}

func (r *tradeRepo) CreateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error {
	return tx.WithContext(ctx).Create(trans).Error
}

func (r *tradeRepo) CreateOrUpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error {
	return tx.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "from_account_id"}, {Name: "nonce"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(trans).Error
}

func (r *tradeRepo) UpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.Transaction) error {
	return tx.WithContext(ctx).Model(&orm.Transaction{}).
		Where(&orm.Transaction{FromAccountID: trans.FromAccountID, Nonce: trans.Nonce}).
		Updates(map[string]interface{}{"status": trans.Status}).Error
}

func (r *tradeRepo) CreateTransactionEvent(ctx context.Context, tx *gorm.DB, event *orm.TransactionEvent) error {
	return tx.WithContext(ctx).Create(event).Error
}

func (r *tradeRepo) GetTransaction(ctx context.Context, tx *gorm.DB, nonce, from int64, status *tcc.Status) (*orm.Transaction, error) {
	var trans orm.Transaction

	q := tx.WithContext(ctx).Where(&orm.Transaction{FromAccountID: from, Nonce: nonce})
	if status != nil {
		q = q.Where("status = ?", *status)
	}

	if err := q.First(&trans).Error; err != nil {
		return nil, err
	}

	return &trans, nil
}
