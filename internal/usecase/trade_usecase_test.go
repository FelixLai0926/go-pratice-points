package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"points/internal/domain"
	"points/internal/domain/command"
	"points/internal/domain/entity"
	"points/internal/domain/repository"
	"points/internal/domain/valueobject"
	"points/internal/infrastructure/distributedlock"
	"points/test/mock"
)

func dummyAccount(userID int64, availableBalance, reservedBalance decimal.Decimal) *entity.Account {
	return &entity.Account{
		UserID:           userID,
		AvailableBalance: valueobject.NewMoneyFromDecimal(availableBalance),
		ReservedBalance:  valueobject.NewMoneyFromDecimal(reservedBalance),
	}
}

func setupTestTradeUsecase(t *testing.T) (
	ctrl *gomock.Controller,
	ctx context.Context,
	mockUow *mock.MockUnitOfWork,
	mockAccRepo *mock.MockAccountRepository,
	mockTxRepo *mock.MockTradeRecordsRepository,
	mockEventRepo *mock.MockTransactionEventRepository,
	mockLocker *mock.MockLocker,
	mockLock *mock.MockLock,
	tradeSvc domain.TradeUsecase,
) {
	ctrl = gomock.NewController(t)
	ctx = context.Background()

	mockUow = mock.NewMockUnitOfWork(ctrl)
	mockAccRepo = mock.NewMockAccountRepository(ctrl)
	mockTxRepo = mock.NewMockTradeRecordsRepository(ctrl)
	mockEventRepo = mock.NewMockTransactionEventRepository(ctrl)
	mockLocker = mock.NewMockLocker(ctrl)
	mockLock = mock.NewMockLock(ctrl)
	mockConfig := mock.NewMockConfig(ctrl)

	//mockTxRepo.EXPECT().GetTradeRecord(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockUow.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(uow repository.UnitOfWork) error) error {
			return fn(mockUow)
		}).AnyTimes()
	mockUow.EXPECT().AccountRepository().Return(mockAccRepo).AnyTimes()
	mockUow.EXPECT().TradeRecordsRepository().Return(mockTxRepo).AnyTimes()
	mockUow.EXPECT().TransactionEventRepository().Return(mockEventRepo).AnyTimes()

	mockConfig.EXPECT().SetDefaultInt("LOCK_DURATION", 5).Return().Times(1)
	mockConfig.EXPECT().SetDefaultInt("RETRY_INTERVAL", 100).Return().Times(1)
	mockConfig.EXPECT().GetInt("LOCK_DURATION").Return(5).Times(1)
	mockConfig.EXPECT().GetInt("RETRY_INTERVAL").Return(100).Times(1)

	tradeSvc = NewTradeUsecase(mockUow, mockLocker, mockConfig)
	return
}

func TestTransfer_AutoConfirmSuccess(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(mockLock, nil).Times(1)
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	mockAccRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
	mockAccRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)

	req := &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 12345,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
		AutoConfirm: true,
	}

	mockAccRepo.EXPECT().ReserveBalance(ctx, req.From, req.Amount).Return(nil).Times(1)
	mockTxRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, nil).Return(nil, nil).Times(1)
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, valueobject.TccPending.Ptr()).Return(&entity.TradeRecords{
		TransactionID: "tx-123",
		FromAccountID: req.From,
		ToAccountID:   req.To,
		Amount:        req.Amount,
		Status:        int32(valueobject.TccPending),
	}, nil).Times(1)

	mockAccRepo.EXPECT().UnreserveBalance(ctx, req.From, req.To, req.Amount).Return(nil).Times(1)
	mockTxRepo.EXPECT().UpdateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)

	err := svc.Transfer(ctx, req)
	if err != nil {
		t.Fatalf("Transfer (auto-confirm) returned error: %v", err)
	}
}

func TestTransfer_NoAutoConfirmSuccess(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	mockAccRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
	mockAccRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)

	req := &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 54321,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(50)),
		AutoConfirm: false,
	}

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, nil).Return(nil, nil).Times(1)
	mockAccRepo.EXPECT().ReserveBalance(ctx, req.From, req.Amount).Return(nil).Times(1)
	mockTxRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)

	err := svc.Transfer(ctx, req)
	if err != nil {
		t.Fatalf("Transfer (no auto-confirm) returned error: %v", err)
	}
}

func TestTransfer_FailureInTransferTransaction(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, _, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	mockAccRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
	mockAccRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)

	req := &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 11111,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(75)),
		AutoConfirm: false,
	}
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, nil).Return(nil, nil).Times(1)
	mockAccRepo.EXPECT().ReserveBalance(ctx, req.From, req.Amount).Return(errors.New("reserve error")).Times(1)

	err := svc.Transfer(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "reserve error") {
		t.Fatalf("Expected reserve error, got: %v", err)
	}
}

func TestTransfer_FailureInConfirm(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	mockAccRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
	mockAccRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(120), decimal.Zero), nil).Times(1)

	req := &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 22222,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(120)),
		AutoConfirm: true,
	}

	mockAccRepo.EXPECT().ReserveBalance(ctx, req.From, req.Amount).Return(nil).Times(1)
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, nil).Return(nil, nil).Times(1)
	mockTxRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(nil, errors.New("get transaction error")).Times(1)

	err := svc.Transfer(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "get transaction error") {
		t.Fatalf("Expected confirm error, got: %v", err)
	}
}

func TestTransfer_LockAcquisitionFailure(t *testing.T) {
	ctrl, ctx, _, _, _, _, mockLocker, _, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("lock acquire failed")).Times(1)

	req := &command.TransferCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 77777,
		},
		Amount:      valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
		AutoConfirm: true,
	}

	err := svc.Transfer(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "lock acquire failed") {
		t.Fatalf("Expected lock acquire error, got: %v", err)
	}
}

func TestManualConfirm_Success(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.ConfirmCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 33333,
		},
	}

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(&entity.TradeRecords{
		TransactionID: "tx-confirm-success",
		FromAccountID: req.From,
		ToAccountID:   req.To,
		Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(80)),
		Status:        int32(valueobject.TccPending),
	}, nil).Times(1)
	mockAccRepo.EXPECT().UnreserveBalance(ctx, req.From, req.To, valueobject.NewMoneyFromDecimal(decimal.NewFromInt(80))).Return(nil).Times(1)
	mockTxRepo.EXPECT().UpdateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)

	err := svc.ManualConfirm(ctx, req)
	if err != nil {
		t.Fatalf("ManualConfirm returned error: %v", err)
	}
}

func TestManualConfirm_Failure(t *testing.T) {
	ctrl, ctx, _, _, mockTxRepo, _, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.ConfirmCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 44444,
		},
	}
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(nil, errors.New("confirm error")).Times(1)

	err := svc.ManualConfirm(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "confirm error") {
		t.Fatalf("Expected confirm error, got: %v", err)
	}
}

func TestManualConfirm_Failure_To_Account_Mismatch(t *testing.T) {
	ctrl, ctx, _, _, mockTxRepo, _, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.ConfirmCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 44444,
		},
	}
	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).
		Return(&entity.TradeRecords{
			TransactionID: "tx-confirm-success",
			FromAccountID: req.From,
			ToAccountID:   int64(9999),
			Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(80)),
			Status:        int32(valueobject.TccPending),
		}, nil).Times(1)

	err := svc.ManualConfirm(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "to account id mismatch") {
		t.Fatalf("Expected cancel error, got: %v", err)
	}
}

func TestCancel_Success(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.CancelCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 55555,
		},
	}

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(&entity.TradeRecords{
		TransactionID: "tx-cancel-success",
		FromAccountID: req.From,
		ToAccountID:   req.To,
		Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(60)),
		Status:        int32(valueobject.TccPending),
	}, nil).Times(1)
	mockAccRepo.EXPECT().UnreserveBalance(ctx, req.From, req.From, valueobject.NewMoneyFromDecimal(decimal.NewFromInt(60))).Return(nil).Times(1)
	mockTxRepo.EXPECT().UpdateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(1)

	err := svc.Cancel(ctx, req)
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
}

func TestCancel_Failure(t *testing.T) {
	ctrl, ctx, _, _, mockTxRepo, _, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.CancelCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 66666,
		},
	}

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(nil, errors.New("cancel error")).Times(1)

	err := svc.Cancel(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "cancel error") {
		t.Fatalf("Expected cancel error, got: %v", err)
	}
}

func TestCancel_Failure_To_Account_Mismatch(t *testing.T) {
	ctrl, ctx, _, _, mockTxRepo, _, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&distributedlock.RedisLock{}, nil).AnyTimes()
	mockLock.EXPECT().Release(ctx).Return(nil).AnyTimes()

	req := &command.CancelCommand{
		BaseCommand: command.BaseCommand{
			From:  1,
			To:    2,
			Nonce: 66666,
		},
	}

	mockTxRepo.EXPECT().GetTradeRecord(ctx, req.Nonce, req.From, gomock.Any()).Return(
		&entity.TradeRecords{
			TransactionID: "tx-cancel-success",
			FromAccountID: req.From,
			ToAccountID:   int64(9999),
			Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(60)),
			Status:        int32(valueobject.TccPending),
		}, nil).Times(1)

	err := svc.Cancel(ctx, req)
	if err == nil || !strings.Contains(err.Error(), "to account id mismatch") {
		t.Fatalf("Expected cancel error, got: %v", err)
	}
}

func TestTransfer_HighConcurrency_ContendLock(t *testing.T) {
	ctrl, ctx, _, mockAccRepo, mockTxRepo, mockEventRepo, mockLocker, mockLock, svc := setupTestTradeUsecase(t)
	defer ctrl.Finish()

	var lock sync.Mutex
	var locked bool

	mockLocker.EXPECT().Acquire(ctx, gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, key string, lockDuration, retryInterval time.Duration) (*distributedlock.RedisLock, error) {
			lock.Lock()
			defer lock.Unlock()
			if !locked {
				locked = true
				return &distributedlock.RedisLock{}, nil
			}
			return nil, fmt.Errorf("lock contention")
		}).AnyTimes()
	mockLock.EXPECT().Release(ctx).DoAndReturn(
		func(ctx context.Context) error {
			lock.Lock()
			defer lock.Unlock()
			time.Sleep(10 * time.Millisecond)
			locked = false
			return nil
		}).AnyTimes()

	transferAmount := decimal.NewFromInt(10)
	mockAccRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, transferAmount, decimal.Zero), nil).AnyTimes()
	mockAccRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).AnyTimes()
	mockAccRepo.EXPECT().ReserveBalance(ctx, int64(1), valueobject.NewMoneyFromDecimal(transferAmount)).AnyTimes().Return(nil)
	mockTxRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).AnyTimes().DoAndReturn(func(ctx context.Context, trans *entity.TradeRecords) error {
		time.Sleep(5 * time.Millisecond)
		trans.TransactionID = uuid.New().String()
		return nil
	})
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).AnyTimes().Return(nil)
	mockTxRepo.EXPECT().GetTradeRecord(gomock.Any(), gomock.Any(), gomock.Any(), nil).Return(nil, nil).AnyTimes()
	mockTxRepo.EXPECT().GetTradeRecord(ctx, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&entity.TradeRecords{
		TransactionID: "dummy-tx",
		FromAccountID: int64(1),
		ToAccountID:   int64(2),
		Amount:        valueobject.NewMoneyFromDecimal(transferAmount),
		Status:        int32(valueobject.TccPending),
	}, nil)
	mockAccRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(2), valueobject.NewMoneyFromDecimal(transferAmount)).AnyTimes().Return(nil)
	mockTxRepo.EXPECT().UpdateTradeRecord(ctx, gomock.Any()).AnyTimes().Return(nil)
	mockEventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).AnyTimes().Return(nil)

	var wg sync.WaitGroup
	var successCount int32
	numTransfers := 1000

	for i := 0; i < numTransfers; i++ {
		wg.Add(1)
		go func(nonce int64) {
			defer wg.Done()
			req := &command.TransferCommand{
				BaseCommand: command.BaseCommand{
					From:  1,
					To:    2,
					Nonce: nonce,
				},
				Amount:      valueobject.NewMoneyFromDecimal(transferAmount),
				AutoConfirm: true,
			}

			err := svc.Transfer(ctx, req)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else if !strings.Contains(err.Error(), "lock") {
				t.Errorf("unexpected error for nonce %d: %v", nonce, err)
			}
		}(int64(1000 + i))
	}
	wg.Wait()

	if successCount != 1 {
		t.Errorf("expected only 1 successful transfer due to lock contention, got %d", successCount)
	}
}
