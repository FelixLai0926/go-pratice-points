package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"points/internal/domain/entity"
	"points/internal/domain/event"
	"points/internal/domain/repository"
	"points/internal/domain/valueobject"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionApplicationService interface {
	TransferTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64, amount valueobject.Money) error
	ConfirmTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64) error
	CancelTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64) error
}

type transactionApplicationService struct{}

func NewTransactionApplicationService() TransactionApplicationService {
	return &transactionApplicationService{}
}

func (ts *transactionApplicationService) TransferTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64, amount valueobject.Money) error {
	toAccount, err := unitOfWork.AccountRepository().GetAccount(ctx, to)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.Wrap(errcode.ErrGetAccount, "transfer phase - get to account", err)
	}

	if toAccount == nil {
		err := unitOfWork.AccountRepository().CreateAccount(ctx, to)
		if err != nil {
			return apperror.Wrap(errcode.ErrCreateAccount, "transfer phase - create account", err)
		}
	}
	tx, err := unitOfWork.TradeRecordsRepository().GetTradeRecord(ctx, nonce, from, nil)
	if tx != nil || (err != nil && !errors.Is(err, gorm.ErrRecordNotFound)) {
		return apperror.Wrap(errcode.ErrConflict, "transfer phase - conflict nonce", err)
	}

	fromAccount, err := unitOfWork.AccountRepository().GetAccount(ctx, from)
	if err != nil {
		return apperror.Wrap(errcode.ErrGetAccount, "transfer phase - get from account", err)
	}

	if err := fromAccount.Reserve(amount); err != nil {
		return err
	}

	if err := unitOfWork.AccountRepository().ReserveBalance(ctx, from, amount); err != nil {
		return apperror.Wrap(errcode.ErrReserveBalance, "transfer phase - reserve balance", err)
	}

	trans := &entity.TradeRecords{
		TransactionID: uuid.New().String(),
		Nonce:         nonce,
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
		Status:        int32(valueobject.TccPending),
	}

	trans.Transfer()

	if err := unitOfWork.TradeRecordsRepository().CreateTradeRecord(ctx, trans); err != nil {
		return apperror.Wrap(errcode.ErrCreateTransaction, "transfer phase - create transaction", err)
	}

	domainEvents := trans.PullEvents()
	if err := ts.saveDomainEvents(ctx, unitOfWork, trans.TransactionID, domainEvents, "transfer phase"); err != nil {
		return err
	}

	return nil
}

func (ts *transactionApplicationService) ConfirmTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64) error {
	trans, err := unitOfWork.TradeRecordsRepository().GetTradeRecord(ctx, nonce, from, valueobject.TccPending.Ptr())
	if err != nil {
		return apperror.Wrap(errcode.ErrGetTransaction, "confirm phase - get transaction", err)
	}
	if trans == nil {
		return apperror.Wrap(errcode.ErrGetTransaction, "confirm phase - get transaction", errors.New("transaction not found"))
	}

	if trans.ToAccountID != to {
		return apperror.Wrap(errcode.ErrInvalidRequest, "confirm phase - to account validation", errors.New("to account id mismatch"))
	}

	if err := unitOfWork.AccountRepository().UnreserveBalance(ctx, from, to, trans.Amount); err != nil {
		return apperror.Wrap(errcode.ErrReserveBalance, "confirm phase - unreserve balance", err)
	}

	trans.Confirm()

	if err := unitOfWork.TradeRecordsRepository().UpdateTradeRecord(ctx, trans); err != nil {
		return apperror.Wrap(errcode.ErrUpdateTransaction, "confirm phase - update transaction", err)
	}

	domainEvents := trans.PullEvents()
	if err := ts.saveDomainEvents(ctx, unitOfWork, trans.TransactionID, domainEvents, "transfer phase"); err != nil {
		return err
	}

	return nil
}

func (ts *transactionApplicationService) CancelTransaction(ctx context.Context, unitOfWork repository.UnitOfWork, nonce, from, to int64) error {
	trans, err := unitOfWork.TradeRecordsRepository().GetTradeRecord(ctx, nonce, from, valueobject.TccPending.Ptr())
	if err != nil {
		return apperror.Wrap(errcode.ErrGetTransaction, "cancel phase - get transaction", err)
	}

	if trans.ToAccountID != to {
		return apperror.Wrap(errcode.ErrInvalidRequest, "cancel phase - to account validation", errors.New("to account id mismatch"))
	}

	if err := unitOfWork.AccountRepository().UnreserveBalance(ctx, from, from, trans.Amount); err != nil {
		return apperror.Wrap(errcode.ErrReserveBalance, "cancel phase - unreserve balance", err)
	}

	trans.Cancel()

	if err := unitOfWork.TradeRecordsRepository().UpdateTradeRecord(ctx, trans); err != nil {
		return err
	}

	domainEvents := trans.PullEvents()
	if err := ts.saveDomainEvents(ctx, unitOfWork, trans.TransactionID, domainEvents, "cancel phase"); err != nil {
		return err
	}

	return nil
}

func (ts *transactionApplicationService) saveDomainEvents(
	ctx context.Context,
	uow repository.UnitOfWork,
	transID string,
	events []event.TransactionEvent,
	phase string,
) error {
	for _, evt := range events {
		payload, err := json.Marshal(evt)
		if err != nil {
			return err
		}

		outboxRecord := entity.TransactionEvent{
			TransactionID: transID,
			EventType:     evt.Action,
			Payload:       string(payload),
		}

		if err := uow.TransactionEventRepository().CreateTransactionEvent(ctx, &outboxRecord); err != nil {
			return apperror.Wrap(errcode.ErrCreateEvent, phase+" - create event", err)
		}
	}
	return nil
}
