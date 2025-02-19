package trade

import (
	"context"

	"github.com/shopspring/decimal"
)

type BaseRequest struct {
	Ctx   context.Context
	From  int64
	To    int64
	Nonce int64
}

type TransferRequest struct {
	BaseRequest
	Amount      decimal.Decimal
	AutoConfirm bool
}

type ConfirmRequest struct {
	BaseRequest
}

type CancelRequest struct {
	BaseRequest
}

type EnsureAccountRequest struct {
	Ctx    context.Context
	UserID int64
}
