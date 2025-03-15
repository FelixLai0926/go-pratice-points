package dto

import "github.com/shopspring/decimal"

type BaseRequest struct {
	From  int64 `json:"from" form:"from" binding:"required"`
	To    int64 `json:"to" form:"to" binding:"required"`
	Nonce int64 `json:"nonce" form:"nonce" binding:"required"`
}

type TransferRequest struct {
	BaseRequest
	Amount      decimal.Decimal `json:"amount" form:"amount" binding:"required"`
	AutoConfirm *bool           `json:"auto_confirm" form:"auto_confirm" default:"true"`
}

type ConfirmRequest struct {
	BaseRequest
}

type CancelRequest struct {
	BaseRequest
}
