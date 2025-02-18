package service

import (
	"points/pkg/models/trade"
)

type TradeService interface {
	Transfer(*trade.TransferRequest) error
	ManualConfirm(*trade.ConfirmRequest) error
	Cancel(*trade.CancelRequest) error
	EnsureDestinationAccount(*trade.EnsureAccountRequest) error
}
