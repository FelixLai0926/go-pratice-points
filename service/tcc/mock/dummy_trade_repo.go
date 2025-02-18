package mock

import (
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"

	"gorm.io/gorm"
)

type DummyTradeRepo struct {
	GetAccountErr     error
	Account           *orm.Account
	GetTransactionErr error
	Transaction       *orm.Transaction
}

func (r *DummyTradeRepo) GetAccount(tx *gorm.DB, userID int32) (*orm.Account, error) {
	if r.GetAccountErr != nil {
		return nil, r.GetAccountErr
	}
	return r.Account, nil
}

func (r *DummyTradeRepo) UpdateAccount(tx *gorm.DB, account *orm.Account) error {
	return nil
}

func (r *DummyTradeRepo) CreateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return nil
}

func (r *DummyTradeRepo) CreateOrUpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return nil
}

func (r *DummyTradeRepo) UpdateTransaction(tx *gorm.DB, trans *orm.Transaction) error {
	return nil
}

func (r *DummyTradeRepo) CreateTransactionEvent(tx *gorm.DB, event *orm.TransactionEvent) error {
	return nil
}

func (r *DummyTradeRepo) GetTransaction(tx *gorm.DB, nonce int64, from int32, status *tcc.Status) (*orm.Transaction, error) {
	if r.GetTransactionErr != nil {
		return nil, r.GetTransactionErr
	}
	return r.Transaction, nil
}
