package mock

import "points/pkg/models/trade"

type DummyTransValidator struct {
	ValidateErr error
}

func (d *DummyTransValidator) ValidateTransferRequest(req *trade.TransferRequest) error {
	return d.ValidateErr
}
func (d *DummyTransValidator) ValidateConfirmRequest(req *trade.ConfirmRequest) error {
	return d.ValidateErr
}

func (d *DummyTransValidator) ValidateCancelRequest(req *trade.CancelRequest) error {
	return d.ValidateErr
}
