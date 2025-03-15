package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/port"
	"points/internal/domain/repository"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/internal/shared/mapper"

	"gorm.io/gorm"
)

var _ repository.TransactionEventRepository = (*transactionEventRepo)(nil)

type transactionEventRepo struct {
	tx     *gorm.DB
	config port.Config
}

func NewTransactionEventRepo(tx *gorm.DB, config port.Config) repository.TransactionEventRepository {
	return &transactionEventRepo{tx: tx, config: config}
}

func (r *transactionEventRepo) CreateTransactionEvent(ctx context.Context, event *entity.TransactionEvent) error {
	ormModel, err := mapper.MapStruct[model.TransactionEvent](r.config, event)
	if err != nil {
		return err
	}
	return r.tx.WithContext(ctx).Create(ormModel).Error
}
