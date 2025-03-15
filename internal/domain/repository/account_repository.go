package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/valueobject"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, userID int64) error
	GetAccount(ctx context.Context, userID int64) (*entity.Account, error)
	ReserveBalance(ctx context.Context, userID int64, amount valueobject.Money) error
	UnreserveBalance(ctx context.Context, from, to int64, amount valueobject.Money) error
}
