package command

import "points/internal/domain/valueobject"

type BaseCommand struct {
	From  int64
	To    int64
	Nonce int64
}

type TransferCommand struct {
	BaseCommand
	Amount      valueobject.Money
	AutoConfirm bool
}

type ConfirmCommand struct {
	BaseCommand
}

type CancelCommand struct {
	BaseCommand
}
