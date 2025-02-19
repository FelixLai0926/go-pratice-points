package mock

import "points/pkg/models/trade"

type DummyTradeService struct {
	TransferErr      error
	ManualConfirmErr error
	CancelErr        error
	Validator        *DummyTransValidator
}

func (d *DummyTradeService) Transfer(rq *trade.TransferRequest) error {
	if err := d.Validator.ValidateTransferRequest(rq); err != nil {
		return err
	}

	return d.TransferErr
}

func (d *DummyTradeService) ManualConfirm(rq *trade.ConfirmRequest) error {
	if err := d.Validator.ValidateConfirmRequest(rq); err != nil {
		return err
	}

	return d.ManualConfirmErr
}

func (d *DummyTradeService) Cancel(rq *trade.CancelRequest) error {
	if err := d.Validator.ValidateCancelRequest(rq); err != nil {
		return err
	}

	return d.CancelErr
}

func (d *DummyTradeService) EnsureDestinationAccount(rq *trade.EnsureAccountRequest) error {
	return nil
}
