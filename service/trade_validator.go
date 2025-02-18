package service

import (
	"points/pkg/models/trade"
)

type TradeValidator interface {
	ValidateTransferRequest(*trade.TransferRequest) error
	ValidateConfirmRequest(*trade.ConfirmRequest) error
	ValidateCancelRequest(*trade.CancelRequest) error
}
