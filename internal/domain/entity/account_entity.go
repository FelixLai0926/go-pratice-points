package entity

import (
	"points/internal/domain/valueobject"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
	"time"
)

type Account struct {
	UserID           int64
	AvailableBalance valueobject.Money
	ReservedBalance  valueobject.Money
	UpdatedAt        time.Time
}

func (a *Account) Reserve(amount valueobject.Money) error {
	if a.AvailableBalance.LessThan(amount) {
		return apperror.Wrap(errcode.ErrInsufficientBalance, "insufficient balance", nil)
	}

	a.AvailableBalance = a.AvailableBalance.Sub(amount)
	a.ReservedBalance = a.ReservedBalance.Add(amount)

	return nil
}

func (a *Account) Unreserve(amount valueobject.Money) error {
	if a.ReservedBalance.LessThan(amount) {
		return apperror.Wrap(errcode.ErrUnreserveBalance, "insufficient balance", nil)
	}
	return nil
}
