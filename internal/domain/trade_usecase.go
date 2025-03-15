package domain

import (
	"context"
	"points/internal/domain/command"
)

type TradeUsecase interface {
	Transfer(ctx context.Context, req *command.TransferCommand) error
	ManualConfirm(ctx context.Context, req *command.ConfirmCommand) error
	Cancel(ctx context.Context, req *command.CancelCommand) error
}
