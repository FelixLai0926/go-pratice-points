package tcc

import (
	stdErrors "errors"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/trade"
	"points/repository"
	"points/service"

	"gorm.io/gorm"
)

type TCCValidator struct {
	db   *gorm.DB
	repo repository.TradeRepository
}

func NewTCCValidator(db *gorm.DB, repo repository.TradeRepository) service.TradeValidator {
	return &TCCValidator{db: db, repo: repo}
}

func (s *TCCValidator) ValidateTransferRequest(rq *trade.TransferRequest) error {
	fromAccount, err := s.repo.GetAccount(s.db, rq.From)
	if err != nil {
		return errors.Wrap(errcode.ErrNotFound, "validate - source account not found", err)
	}

	if fromAccount.AvailableBalance < rq.Amount {
		return errors.Wrap(errcode.ErrInsufficientBalance, "validate - insufficient balance", nil)
	}

	existingTrans, err := s.repo.GetTransaction(s.db, rq.Nonce, rq.From, nil)
	if err == nil && existingTrans != nil {
		return errors.Wrap(errcode.ErrConflict, "validate - nonce already used", err)
	}

	return nil
}

func (s *TCCValidator) ValidateConfirmRequest(rq *trade.ConfirmRequest) error {
	if err := s.validateAccountAndTransaction(rq.From, rq.Nonce); err != nil {
		return err
	}

	return nil
}

func (s *TCCValidator) ValidateCancelRequest(rq *trade.CancelRequest) error {
	if err := s.validateAccountAndTransaction(rq.From, rq.Nonce); err != nil {
		return err
	}

	return nil
}

func (s *TCCValidator) validateAccountAndTransaction(from int32, nonce int64) error {
	fromAccount, err := s.repo.GetAccount(s.db, from)
	if err != nil {
		return errors.Wrap(errcode.ErrNotFound, "validate - source account not found", err)
	}

	pendingStatus := tcc.Pending
	pendingTransaction, err := s.repo.GetTransaction(s.db, nonce, from, &pendingStatus)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(errcode.ErrGetTransaction, "validate - transaction not found", err)
		}
		return errors.Wrap(errcode.ErrGetTransaction, "validate - error retrieving transaction", err)
	}
	if pendingTransaction == nil {
		return errors.Wrap(errcode.ErrGetTransaction, "validate - transaction not found", nil)
	}

	if fromAccount.ReservedBalance < pendingTransaction.Amount {
		return errors.Wrap(errcode.ErrInsufficientBalance, "validate - insufficient balance", nil)
	}

	return nil
}
