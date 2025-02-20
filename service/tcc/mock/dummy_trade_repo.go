package mock

import (
	"context"
	"errors"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"
	"time"

	"gorm.io/gorm"
)

type DummyTradeRepo struct {
	FailUpdateAccount     bool
	FailCreateTransaction bool
	Delay                 *time.Duration
	DummyAccount          *orm.Account
	DummyTransaction      *orm.TransactionDAO
}

func (r *DummyTradeRepo) CreateAccount(ctx context.Context, tx *gorm.DB, userID int64) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	return nil
}
func (r *DummyTradeRepo) GetAccount(ctx context.Context, tx *gorm.DB, userID int64) (*orm.Account, error) {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	return r.DummyAccount, nil
}
func (r *DummyTradeRepo) UpdateAccount(ctx context.Context, tx *gorm.DB, account *orm.Account) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	if r.FailUpdateAccount {
		return errors.New("simulated update account failure")
	}
	return nil
}
func (r *DummyTradeRepo) CreateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.TransactionDAO) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	if r.FailCreateTransaction {
		return errors.New("simulated create transaction failure")
	}
	return nil
}
func (r *DummyTradeRepo) CreateOrUpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.TransactionDAO) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	if r.FailCreateTransaction {
		return errors.New("simulated create transaction failure")
	}
	return nil
}
func (r *DummyTradeRepo) UpdateTransaction(ctx context.Context, tx *gorm.DB, trans *orm.TransactionDAO) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	return nil
}
func (r *DummyTradeRepo) CreateTransactionEvent(ctx context.Context, tx *gorm.DB, event *orm.Transaction_event) error {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	return nil
}
func (r *DummyTradeRepo) GetTransaction(ctx context.Context, tx *gorm.DB, nonce, from int64, status *tcc.Status) (*orm.TransactionDAO, error) {
	if r.Delay != nil {
		time.Sleep(*r.Delay)
	}
	return r.DummyTransaction, nil
}
