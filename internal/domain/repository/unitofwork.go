package repository

import "context"

type UnitOfWork interface {
	AccountRepository() AccountRepository
	TradeRecordsRepository() TradeRecordsRepository
	TransactionEventRepository() TransactionEventRepository
	Transaction(context.Context, func(UnitOfWork) error) error
}
