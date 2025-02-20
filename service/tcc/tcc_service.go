package tcc

import (
	stdErrors "errors"
	"fmt"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/enum/tcc"
	event "points/pkg/models/eventpayload"
	"points/pkg/models/orm"
	"points/pkg/models/trade"
	"points/pkg/module/config"
	"points/pkg/module/distributedlock"
	"points/repository"
	"points/service"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TCCService struct {
	db         *gorm.DB
	repo       repository.TradeRepository
	validator  service.TradeValidator
	lockClient distributedlock.LockClient
}

func NewTCCService(db *gorm.DB, repo repository.TradeRepository, validator service.TradeValidator, lockClient distributedlock.LockClient) service.TradeService {
	return &TCCService{db: db, repo: repo, validator: validator, lockClient: lockClient}
}
func (s *TCCService) Transfer(rq *trade.TransferRequest) error {
	if err := s.validator.ValidateTransferRequest(rq); err != nil {
		return err
	}
	return s.internalTransfer(rq)
}

func (s *TCCService) ManualConfirm(rq *trade.ConfirmRequest) error {
	if err := s.validator.ValidateConfirmRequest(rq); err != nil {
		return err
	}
	return s.internalManualConfirm(rq)
}

func (s *TCCService) Cancel(rq *trade.CancelRequest) error {
	if err := s.validator.ValidateCancelRequest(rq); err != nil {
		return err
	}
	return s.internalCancel(rq)
}

func (s *TCCService) EnsureDestinationAccount(rq *trade.EnsureAccountRequest) error {
	return s.db.WithContext(rq.Ctx).Transaction(func(tx *gorm.DB) error {
		account, err := s.repo.GetAccount(rq.Ctx, tx, rq.UserID)
		if err != nil && !stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(errcode.ErrGetAccount, "create destination account - get account", err)
		}

		if account != nil {
			return nil
		}
		return s.repo.CreateAccount(rq.Ctx, tx, rq.UserID)
	})
}
func (s *TCCService) internalTransfer(rq *trade.TransferRequest) error {
	lockKey := getLockKey(rq.From, rq.To)
	ttl := getLockTTL()

	return distributedlock.WithLock(rq.Ctx, s.lockClient, lockKey, ttl, func() error {
		return s.db.WithContext(rq.Ctx).Transaction(func(tx *gorm.DB) error {
			fromAccount, err := s.repo.GetAccount(rq.Ctx, tx, rq.From)
			if err != nil {
				return errors.Wrap(errcode.ErrGetAccount, "try phase - get account", err)
			}

			fromAccount.AvailableBalance = fromAccount.AvailableBalance.Sub(rq.Amount)
			fromAccount.ReservedBalance = fromAccount.ReservedBalance.Add(rq.Amount)
			if err := s.updateAccountAndCreateTransaction(tx, &rq.BaseRequest, fromAccount, rq.Amount, tcc.Pending); err != nil {
				return err
			}

			if !rq.AutoConfirm {
				return nil
			}

			if err := s.confirm(tx, &rq.BaseRequest, fromAccount); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *TCCService) internalManualConfirm(rq *trade.ConfirmRequest) error {
	lockKey := getLockKey(rq.From, rq.To)
	ttl := getLockTTL()

	return distributedlock.WithLock(rq.Ctx, s.lockClient, lockKey, ttl, func() error {
		return s.db.WithContext(rq.Ctx).Transaction(func(tx *gorm.DB) error {
			fromAccount, err := s.repo.GetAccount(rq.Ctx, tx, rq.From)
			if err != nil {
				return errors.Wrap(errcode.ErrGetAccount, "confirm phase - get account", err)
			}

			if err := s.confirm(tx, &rq.BaseRequest, fromAccount); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *TCCService) internalCancel(rq *trade.CancelRequest) error {
	lockKey := getLockKey(rq.From, rq.To)
	ttl := getLockTTL()

	return distributedlock.WithLock(rq.Ctx, s.lockClient, lockKey, ttl, func() error {
		return s.db.WithContext(rq.Ctx).Transaction(func(tx *gorm.DB) error {
			fromAccount, err := s.repo.GetAccount(rq.Ctx, tx, rq.From)
			if err != nil {
				return errors.Wrap(errcode.ErrGetAccount, "cancel phase - get source account", err)
			}

			pendingStatus := tcc.Pending
			trans, err := s.repo.GetTransaction(rq.Ctx, tx, rq.Nonce, rq.From, &pendingStatus)
			if err != nil {
				return errors.Wrap(errcode.ErrGetTransaction, "cancel phase - get transaction", err)
			}

			fromAccount.ReservedBalance = fromAccount.ReservedBalance.Sub(trans.Amount)
			fromAccount.AvailableBalance = fromAccount.AvailableBalance.Add(trans.Amount)

			if err := s.updateAccountAndCreateTransaction(tx, &rq.BaseRequest, fromAccount, trans.Amount, tcc.Canceled); err != nil {
				return err
			}

			return nil
		})
	})
}

func (s *TCCService) confirm(tx *gorm.DB, rq *trade.BaseRequest, fromAccount *orm.Account) error {
	pendingStatus := tcc.Pending
	trans, err := s.repo.GetTransaction(rq.Ctx, tx, rq.Nonce, rq.From, &pendingStatus)
	if err != nil {
		return errors.Wrap(errcode.ErrGetTransaction, "confirm phase - get transaction", err)
	}

	fromAccount.ReservedBalance = fromAccount.ReservedBalance.Sub(trans.Amount)
	if err := s.repo.UpdateAccount(rq.Ctx, tx, fromAccount); err != nil {
		return errors.Wrap(errcode.ErrUpdateAccount, "confirm phase - update account", err)
	}

	toAccount, err := s.repo.GetAccount(rq.Ctx, tx, rq.To)
	if err != nil {
		return errors.Wrap(errcode.ErrGetAccount, "confirm phase - get destination account", err)
	}

	toAccount.AvailableBalance = toAccount.AvailableBalance.Add(trans.Amount)
	if err := s.updateAccountAndCreateTransaction(tx, rq, toAccount, trans.Amount, tcc.Confirmed); err != nil {
		return err
	}

	return nil
}

func (s *TCCService) updateAccountAndCreateTransaction(tx *gorm.DB, rq *trade.BaseRequest, account *orm.Account, amount decimal.Decimal, status tcc.Status) error {
	if err := s.repo.UpdateAccount(rq.Ctx, tx, account); err != nil {
		return errors.Wrap(errcode.ErrUpdateAccount, fmt.Sprintf("%s phase - update account", status.String()), err)
	}

	trans := &orm.TransactionDAO{
		TransactionID: uuid.New().String(),
		Nonce:         rq.Nonce,
		FromAccountID: rq.From,
		ToAccountID:   rq.To,
		Amount:        amount,
		Status:        int32(status),
	}

	if err := s.repo.CreateOrUpdateTransaction(rq.Ctx, tx, trans); err != nil {
		return errors.Wrap(errcode.ErrCreateTransaction, fmt.Sprintf("%s phase - create transaction", status.String()), err)
	}

	tryEventPayload := event.TransferPayload{
		Action: status.String(),
		Amount: amount,
	}

	tryEventPayloadJSON, err := tryEventPayload.ToJSON()
	if err != nil {
		return errors.Wrap(errcode.ErrPayloadMarshal, fmt.Sprintf("%s phase - failed to marshal payload", status.String()), err)
	}

	tryEvent := &orm.Transaction_event{
		TransactionID: trans.TransactionID,
		EventType:     status.String(),
		Payload:       tryEventPayloadJSON,
	}

	if err := s.repo.CreateTransactionEvent(rq.Ctx, tx, tryEvent); err != nil {
		return errors.Wrap(errcode.ErrCreateEvent, fmt.Sprintf("%s phase - create transaction event", status.String()), err)
	}

	return nil
}

func getLockKey(from, to int64) string {
	if from < to {
		return fmt.Sprintf("transfer_lock:%d:%d", from, to)
	}
	return fmt.Sprintf("transfer_lock:%d:%d", to, from)
}

func getLockTTL() time.Duration {
	ttlSeconds, err := config.GetEnvOrDefault("LOCK_TTL", 5, strconv.Atoi)
	if err != nil {
		zap.L().Error("get env error, use default value", zap.Int("defaultValue", 5), zap.Error(err))
		ttlSeconds = 5
	}
	return time.Duration(ttlSeconds) * time.Second
}
