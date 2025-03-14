package repository

import (
	"context"
	"points/internal/domain/entity"
)

type TransactionEventRepository interface {
	CreateTransactionEvent(ctx context.Context, event *entity.TransactionEvent) error
}
