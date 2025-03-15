package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/valueobject"
)

type TradeRecordsRepository interface {
	CreateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error
	CreateOrUpdateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error
	UpdateTradeRecord(ctx context.Context, trans *entity.TradeRecords) error
	GetTradeRecord(ctx context.Context, nonce, from int64, status *valueobject.TccStatus) (*entity.TradeRecords, error)
}
