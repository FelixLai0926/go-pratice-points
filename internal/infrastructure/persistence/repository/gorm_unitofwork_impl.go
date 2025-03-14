package repository

import (
	"context"
	"points/internal/domain/port"
	"points/internal/domain/repository"

	"gorm.io/gorm"
)

type gormUnitOfWorkImpl struct {
	db                *gorm.DB
	tx                *gorm.DB
	isTransaction     bool
	accountRepository repository.AccountRepository
	transactionRepo   repository.TradeRecordsRepository
	eventRepository   repository.TransactionEventRepository
	config            port.Config
}

var _ repository.UnitOfWork = (*gormUnitOfWorkImpl)(nil)

func NewGormUnitOfWorkImpl(db *gorm.DB, config port.Config) repository.UnitOfWork {
	return &gormUnitOfWorkImpl{
		db:                db,
		isTransaction:     false,
		config:            config,
		accountRepository: nil,
		transactionRepo:   nil,
		eventRepository:   nil,
	}
}

func (u *gormUnitOfWorkImpl) AccountRepository() repository.AccountRepository {
	if u.accountRepository == nil {
		u.accountRepository = NewAccountRepo(u.getCurrentDB(), u.config)
	}
	return u.accountRepository
}

func (u *gormUnitOfWorkImpl) TradeRecordsRepository() repository.TradeRecordsRepository {
	if u.transactionRepo == nil {
		u.transactionRepo = NewTradeRecordsRepo(u.getCurrentDB(), u.config)
	}
	return u.transactionRepo
}

func (u *gormUnitOfWorkImpl) TransactionEventRepository() repository.TransactionEventRepository {
	if u.eventRepository == nil {
		u.eventRepository = NewTransactionEventRepo(u.getCurrentDB(), u.config)
	}
	return u.eventRepository
}

func (u *gormUnitOfWorkImpl) Transaction(ctx context.Context, fn func(uow repository.UnitOfWork) error) error {
	if u.isTransaction {
		return fn(u)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		uow := &gormUnitOfWorkImpl{
			tx:                tx,
			isTransaction:     true,
			accountRepository: nil,
			transactionRepo:   nil,
			eventRepository:   nil,
			config:            u.config,
		}
		return fn(uow)
	})
}

func (u *gormUnitOfWorkImpl) DB(ctx context.Context) *gorm.DB {
	return u.getCurrentDB().WithContext(ctx)
}

func (u *gormUnitOfWorkImpl) getCurrentDB() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}
	return u.db
}
