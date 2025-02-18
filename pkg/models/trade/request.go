package trade

import "context"

type BaseRequest struct {
	Ctx   context.Context
	From  int32
	To    int32
	Nonce int64
}

type TransferRequest struct {
	BaseRequest
	Amount      float64
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
	UserID int32
}
