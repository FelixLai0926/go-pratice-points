package di

import (
	"points/internal/domain"
	"points/internal/domain/port"
	"points/internal/domain/repository"
	"points/internal/usecase"
	"points/internal/usecase/locking"
	"points/internal/usecase/transaction"

	"go.uber.org/fx"
)

var ApplicationModule = fx.Options(
	fx.Provide(func(config port.Config, locker domain.Locker) locking.AccountLockApplicationService {
		return locking.NewAccountLockService(locker, config)
	}),
	fx.Provide(transaction.NewTransactionApplicationService),
	fx.Provide(func(uow repository.UnitOfWork, locker domain.Locker, config port.Config) domain.TradeUsecase {
		return usecase.NewTradeUsecase(uow, locker, config)
	}),
)
