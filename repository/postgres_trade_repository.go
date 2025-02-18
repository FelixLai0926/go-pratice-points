package repository

import (
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tradeRepo struct{}

func NewTradeRepo() TradeRepository {
	return &tradeRepo{}
}

func (r *tradeRepo) CreateAccount(tx *gorm.DB, userID int32) error {
	return tx.Create(&orm.Account{
		UserID:           userID,
		AvailableBalance: 0,
		ReservedBalance:  0,
	}).Error
}

func (r *tradeRepo) GetAccount(tx *gorm.DB, userID int32) (*orm.Account, error) {
	var account orm.Account
	if err := tx.First(&account, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *tradeRepo) UpdateAccount(tx *gorm.DB, account *orm.Account) error {
	return tx.Save(account).Error
}

func (r *tradeRepo) CreateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return tx.Create(trans).Error
}

func (r *tradeRepo) CreateOrUpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "from_account_id"}, {Name: "nonce"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(trans).Error
}

func (r *tradeRepo) UpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return tx.Model(&orm.Transaction{}).
		Where("from_account_id = ? AND nonce = ?", trans.FromAccountID, trans.Nonce).
		Updates(trans).Error
}

func (r *tradeRepo) CreateTransactionEvent(tx *gorm.DB, event *orm.TransactionEvent) error {
	return tx.Create(event).Error
}

func (r *tradeRepo) GetTransaction(tx *gorm.DB, nonce int64, from int32, status *tcc.Status) (*orm.Transaction, error) {
	var trans orm.Transaction
	q := tx.Where("nonce = ? AND from_account_id = ?", nonce, from)
	if status != nil {
		q = q.Where("status = ?", *status)
	}

	if err := q.First(&trans).Error; err != nil {
		return nil, err
	}

	return &trans, nil
}
