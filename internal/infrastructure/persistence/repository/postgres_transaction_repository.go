package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/port"
	"points/internal/domain/repository"
	"points/internal/domain/valueobject"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/internal/shared/mapper"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repository.TradeRecordsRepository = (*tradeRecordsRepo)(nil)

type tradeRecordsRepo struct {
	tx     *gorm.DB
	config port.Config
}

func NewTradeRecordsRepo(tx *gorm.DB, config port.Config) repository.TradeRecordsRepository {
	return &tradeRecordsRepo{tx: tx, config: config}
}

func (r *tradeRecordsRepo) CreateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error {
	ormModel, err := mapper.MapStruct[model.TradeRecord](r.config, trans)
	if err != nil {
		return err
	}

	return r.tx.WithContext(ctx).Create(ormModel).Error
}

func (r *tradeRecordsRepo) CreateOrUpdateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error {
	ormModel, err := mapper.MapStruct[model.TradeRecord](r.config, trans)
	if err != nil {
		return err
	}
	return r.tx.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "from_account_id"}, {Name: "nonce"}},
		DoUpdates: clause.AssignmentColumns([]string{"status"}),
	}).Create(ormModel).Error
}

func (r *tradeRecordsRepo) UpdateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error {
	return r.tx.WithContext(ctx).Model(&model.TradeRecord{}).
		Where(&model.TradeRecord{
			TransactionID: trans.TransactionID,
			FromAccountID: trans.FromAccountID,
			Nonce:         trans.Nonce,
		}).
		Updates(map[string]interface{}{"status": trans.Status}).Error
}

func (r *tradeRecordsRepo) GetTradeRecord(ctx context.Context, nonce, from int64, status *valueobject.TccStatus) (*entity.TradeRecords, error) {
	var trans model.TradeRecord

	q := r.tx.WithContext(ctx).Where(&model.TradeRecord{FromAccountID: from, Nonce: nonce})
	if status != nil {
		q = q.Where("status = ?", *status)
	}

	if err := q.First(&trans).Error; err != nil {
		return nil, err
	}

	domainModel, err := mapper.MapStruct[entity.TradeRecords](r.config, &trans)
	if err != nil {
		return nil, err
	}
	return domainModel, nil
}
