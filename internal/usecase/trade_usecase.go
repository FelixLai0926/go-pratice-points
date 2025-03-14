package usecase

import (
	"context"
	"points/internal/domain"
	"points/internal/domain/command"
	"points/internal/domain/port"
	"points/internal/domain/repository"
	"points/internal/usecase/locking"
	"points/internal/usecase/transaction"
)

type tradeUsecase struct {
	unitOfWork         repository.UnitOfWork
	lockService        locking.AccountLockApplicationService
	transactionService transaction.TransactionApplicationService
}

func NewTradeUsecase(unitOfWork repository.UnitOfWork, locker domain.Locker, config port.Config) domain.TradeUsecase {
	return &tradeUsecase{
		unitOfWork:         unitOfWork,
		lockService:        locking.NewAccountLockService(locker, config),
		transactionService: transaction.NewTransactionApplicationService(),
	}
}

func (s *tradeUsecase) Transfer(ctx context.Context, req *command.TransferCommand) error {
	return s.lockService.WithAccountTradeLock(ctx, req.From, req.To, func() error {
		return s.unitOfWork.Transaction(ctx, func(u repository.UnitOfWork) error {
			if err := s.transactionService.TransferTransaction(ctx, u, req.Nonce, req.From, req.To, req.Amount); err != nil {
				return err
			}

			if !req.AutoConfirm {
				return nil
			}

			if err := s.confirm(ctx, &req.BaseCommand, u); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *tradeUsecase) ManualConfirm(ctx context.Context, req *command.ConfirmCommand) error {
	return s.lockService.WithAccountTradeLock(ctx, req.From, req.To, func() error {
		return s.unitOfWork.Transaction(ctx, func(u repository.UnitOfWork) error {
			if err := s.confirm(ctx, &req.BaseCommand, u); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *tradeUsecase) Cancel(ctx context.Context, req *command.CancelCommand) error {
	return s.lockService.WithAccountTradeLock(ctx, req.From, req.To, func() error {
		return s.unitOfWork.Transaction(ctx, func(u repository.UnitOfWork) error {
			if err := s.transactionService.CancelTransaction(ctx, u, req.Nonce, req.From, req.To); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *tradeUsecase) confirm(ctx context.Context, rq *command.BaseCommand, unitOfWork repository.UnitOfWork) error {
	err := s.transactionService.ConfirmTransaction(ctx, unitOfWork, rq.Nonce, rq.From, rq.To)
	if err != nil {
		return err
	}
	return nil
}
