package mock

import "points/pkg/models/trade"

type DummyTradeService struct {
	TransferErr      error
	ManualConfirmErr error
	CancelErr        error
}

func (d *DummyTradeService) Transfer(rq *trade.TransferRequest) error {
	return d.TransferErr
}

func (d *DummyTradeService) ManualConfirm(rq *trade.ConfirmRequest) error {
	return d.ManualConfirmErr
}

func (d *DummyTradeService) Cancel(rq *trade.CancelRequest) error {
	return d.CancelErr
}

func (d *DummyTradeService) EnsureDestinationAccount(rq *trade.EnsureAccountRequest) error {
	return nil
}
