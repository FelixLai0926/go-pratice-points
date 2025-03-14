package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/port"
	"points/internal/domain/repository"
	"points/internal/domain/valueobject"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/internal/shared/mapper"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var _ repository.AccountRepository = (*accountRepo)(nil)

type accountRepo struct {
	tx     *gorm.DB
	config port.Config
}

func NewAccountRepo(tx *gorm.DB, config port.Config) repository.AccountRepository {
	return &accountRepo{tx: tx, config: config}
}

func (r *accountRepo) CreateAccount(ctx context.Context, userID int64) error {
	return r.tx.WithContext(ctx).Create(&model.Account{
		UserID:           userID,
		AvailableBalance: decimal.Zero,
		ReservedBalance:  decimal.Zero,
	}).Error
}

func (r *accountRepo) GetAccount(ctx context.Context, userID int64) (*entity.Account, error) {
	var account model.Account
	err := r.tx.WithContext(ctx).
		Where(&model.Account{UserID: userID}).
		First(&account).Error
	if err != nil {
		return nil, err
	}

	domainAccount, err := mapper.MapStruct[entity.Account](r.config, &account)
	if err != nil {
		return nil, err
	}

	return domainAccount, nil
}

func (r *accountRepo) ReserveBalance(ctx context.Context, userID int64, amount valueobject.Money) error {
	return r.tx.WithContext(ctx).Model(&model.Account{}).
		Where(&model.Account{UserID: userID}).
		Updates(map[string]interface{}{
			"available_balance": gorm.Expr("available_balance - ?", amount.Value()),
			"reserved_balance":  gorm.Expr("reserved_balance + ?", amount.Value()),
		}).Error
}

func (r *accountRepo) UnreserveBalance(ctx context.Context, from, to int64, amount valueobject.Money) error {
	err := r.tx.WithContext(ctx).Model(&model.Account{}).
		Where(&model.Account{UserID: from}).
		Updates(map[string]interface{}{
			"reserved_balance": gorm.Expr("reserved_balance - ?", amount.Value()),
		}).Error

	if err != nil {
		return err
	}
	err = r.tx.WithContext(ctx).Model(&model.Account{}).
		Where(&model.Account{UserID: to}).
		Updates(map[string]interface{}{
			"available_balance": gorm.Expr("available_balance + (?::numeric)", amount.Value()),
		}).Error

	return err
}
