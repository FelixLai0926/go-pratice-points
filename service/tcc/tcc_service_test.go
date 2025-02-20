package tcc

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"
	"points/pkg/models/trade"
	"points/pkg/module/distributedlock"
	"points/pkg/module/test"
	"points/repository"
	"points/service"
	"points/service/tcc/mock"

	"github.com/alicebob/miniredis/v2"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type testCase[T any] struct {
	name                 string
	setup                func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func())
	request              *T
	expectedErrSubstring string
	validate             func(t *testing.T, db *gorm.DB)
}

func TestTransfer(t *testing.T) {
	testCases := []testCase[trade.TransferRequest]{
		{
			name: "AutoConfirm success",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 12345,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: true,
			},
			expectedErrSubstring: "",
			validate: func(t *testing.T, db *gorm.DB) {
				var fromAcc orm.Account
				if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query from account: %v", err)
				}
				if !fromAcc.AvailableBalance.Equal(decimal.NewFromInt(900)) {
					t.Errorf("expected from account available balance 900, got %s", fromAcc.AvailableBalance.String())
				}
				var toAcc orm.Account
				if err := db.First(&toAcc, "user_id = ?", 2).Error; err != nil {
					t.Fatalf("failed to query to account: %v", err)
				}
				if !toAcc.AvailableBalance.Equal(decimal.NewFromInt(600)) {
					t.Errorf("expected to account available balance 600, got %s", toAcc.AvailableBalance.String())
				}
			},
		},
		{
			name: "NonAutoConfirm success",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 55555,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: false,
			},
			expectedErrSubstring: "",
			validate: func(t *testing.T, db *gorm.DB) {
				var fromAcc orm.Account
				if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query from account: %v", err)
				}
				if !fromAcc.AvailableBalance.Equal(decimal.NewFromInt(900)) {
					t.Errorf("expected from account available balance 900, got %s", fromAcc.AvailableBalance.String())
				}
				if !fromAcc.ReservedBalance.Equal(decimal.NewFromInt(100)) {
					t.Errorf("expected from account reserved balance 100, got %s", fromAcc.ReservedBalance.String())
				}
				var toAcc orm.Account
				if err := db.First(&toAcc, "user_id = ?", 2).Error; err != nil {
					t.Fatalf("failed to query to account: %v", err)
				}
				if !toAcc.AvailableBalance.Equal(decimal.NewFromInt(500)) {
					t.Errorf("expected to account available balance 500, got %s", toAcc.AvailableBalance.String())
				}
			},
		},
		{
			name: "UpdateAccount failure",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				dummyRepo := &mock.DummyTradeRepo{
					FailUpdateAccount: true,
					DummyAccount: &orm.Account{
						UserID:           1,
						AvailableBalance: decimal.NewFromInt(1000),
						ReservedBalance:  decimal.Zero,
					},
					DummyTransaction: nil,
				}
				db, _, svc, miniredis := setupMockTestService(t, dummyRepo)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 70000,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: true,
			},
			expectedErrSubstring: "simulated update account failure",
			validate:             nil,
		},
		{
			name: "CreateTransaction failure",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				dummyRepo := &mock.DummyTradeRepo{
					FailCreateTransaction: true,
					DummyAccount: &orm.Account{
						UserID:           1,
						AvailableBalance: decimal.NewFromInt(1000),
						ReservedBalance:  decimal.Zero,
					},
					DummyTransaction: nil,
				}
				db, _, svc, miniredis := setupMockTestService(t, dummyRepo)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 70001,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: true,
			},
			expectedErrSubstring: "simulated create transaction failure",
			validate:             nil,
		},
		{
			name: "Lock failure",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db = test.NewTestContainerDB(t)
				test.SetupAccounts(t, db)
				repo := repository.NewTradeRepo()
				validator := NewTCCValidator(db, repo)
				failingLockClient := distributedlock.NewFailingLockClient()
				svc = NewTCCService(db, repo, validator, failingLockClient)
				return db, svc, func() {}
			},
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 888,
				},
				Amount:      decimal.NewFromInt(50),
				AutoConfirm: true,
			},
			expectedErrSubstring: "lock",
			validate:             nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, svc, cleanup := tc.setup(t)
			if cleanup != nil {
				defer cleanup()
			}
			err := svc.Transfer(tc.request)
			if tc.expectedErrSubstring != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedErrSubstring) {
					t.Fatalf("expected error containing %q, got: %v", tc.expectedErrSubstring, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
				if tc.validate != nil {
					tc.validate(t, db)
				}
			}
		})
	}
}

func TestConfirm(t *testing.T) {
	testCases := []testCase[trade.ConfirmRequest]{
		{
			name: "Confirm success",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				pendingTx := &orm.TransactionDAO{
					TransactionID: "tx-confirm-success",
					Nonce:         456,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(100),
					Status:        int32(tcc.Pending),
				}
				if err := db.Create(pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
				var acc orm.Account
				if err := db.First(&acc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to get account: %v", err)
				}
				acc.ReservedBalance = decimal.NewFromInt(100)
				if err := db.Save(&acc).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 456,
				},
			},
			expectedErrSubstring: "",
			validate: func(t *testing.T, db *gorm.DB) {
				var tx orm.TransactionDAO
				if err := db.Where("nonce = ? AND from_account_id = ?", 456, 1).First(&tx).Error; err != nil {
					t.Fatalf("failed to query transaction: %v", err)
				}
				if tx.Status != int32(tcc.Confirmed) {
					t.Errorf("expected transaction status Confirmed, got %d", tx.Status)
				}
				var fromAcc orm.Account
				if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query from account: %v", err)
				}
				if !fromAcc.ReservedBalance.Equal(decimal.Zero) {
					t.Errorf("expected from account reserved balance 0, got %s", fromAcc.ReservedBalance.String())
				}
				var toAcc orm.Account
				if err := db.First(&toAcc, "user_id = ?", 2).Error; err != nil {
					t.Fatalf("failed to query to account: %v", err)
				}
				if !toAcc.AvailableBalance.Equal(decimal.NewFromInt(600)) {
					t.Errorf("expected to account available balance 600, got %s", toAcc.AvailableBalance.String())
				}
			},
		},
		{
			name: "Confirm fails due to missing transaction",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 999,
				},
			},
			expectedErrSubstring: "record not found",
			validate:             nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, svc, cleanup := tc.setup(t)
			if cleanup != nil {
				defer cleanup()
			}
			err := svc.ManualConfirm(tc.request)
			if tc.expectedErrSubstring != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedErrSubstring) {
					t.Fatalf("expected error containing %q, got: %v", tc.expectedErrSubstring, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
				if tc.validate != nil {
					tc.validate(t, db)
				}
			}
		})
	}
}

func TestCancel(t *testing.T) {
	testCases := []testCase[trade.CancelRequest]{
		{
			name: "Cancel success",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				pendingTx := &orm.TransactionDAO{
					TransactionID: "tx-cancel-success",
					Nonce:         789,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(50),
					Status:        int32(tcc.Pending),
				}
				if err := db.Create(pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
				var acc orm.Account
				if err := db.First(&acc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query account: %v", err)
				}
				acc.AvailableBalance = decimal.NewFromInt(950)
				acc.ReservedBalance = decimal.NewFromInt(50)
				if err := db.Save(&acc).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 789,
				},
			},
			expectedErrSubstring: "",
			validate: func(t *testing.T, db *gorm.DB) {
				var fromAcc orm.Account
				if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query from account: %v", err)
				}
				if !fromAcc.ReservedBalance.Equal(decimal.Zero) {
					t.Errorf("expected reserved balance 0 after cancel, got %s", fromAcc.ReservedBalance.String())
				}
				if !fromAcc.AvailableBalance.Equal(decimal.NewFromInt(1000)) {
					t.Errorf("expected available balance 1000 after cancel, got %s", fromAcc.AvailableBalance.String())
				}
			},
		},
		{
			name: "Cancel failure - record not found",
			setup: func(t *testing.T) (db *gorm.DB, svc service.TradeService, cleanup func()) {
				db, _, svc, miniredis := setupTestService(t)
				return db, svc, func() { miniredis.Close() }
			},
			request: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 999,
				},
			},
			expectedErrSubstring: "record not found",
			validate:             nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, svc, cleanup := tc.setup(t)
			if cleanup != nil {
				defer cleanup()
			}
			err := svc.Cancel(tc.request)
			if tc.expectedErrSubstring != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedErrSubstring) {
					t.Fatalf("expected error containing %q, got: %v", tc.expectedErrSubstring, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
				if tc.validate != nil {
					tc.validate(t, db)
				}
			}
		})
	}
}

func TestTransfer_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	delay := time.Duration(10 * time.Millisecond)
	dummyRepo := &mock.DummyTradeRepo{
		DummyAccount: &orm.Account{
			UserID:           1,
			AvailableBalance: decimal.NewFromInt(1000),
			ReservedBalance:  decimal.Zero,
		},
		DummyTransaction: nil,
		Delay:            &delay,
	}

	_, _, svc, miniredis := setupMockTestService(t, dummyRepo)
	defer miniredis.Close()

	req := &trade.TransferRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   ctx,
			From:  1,
			To:    2,
			Nonce: 1,
		},
		Amount:      decimal.NewFromInt(100),
		AutoConfirm: true,
	}

	err := svc.Transfer(req)
	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("expected context deadline exceeded error, got: %v", err)
	}
}

func TestTransfer_HighConcurrency(t *testing.T) {
	db, _, svc, miniredis := setupTestService(t)
	defer miniredis.Close()

	numTransfers := 100
	transferAmount := decimal.NewFromInt(10)
	var wg sync.WaitGroup
	var lockFailures int32

	for i := 0; i < numTransfers; i++ {
		wg.Add(1)
		go func(nonce int64) {
			defer wg.Done()
			req := &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: nonce,
				},
				Amount:      transferAmount,
				AutoConfirm: true,
			}
			if err := svc.Transfer(req); err != nil {
				if strings.Contains(err.Error(), "lock") {
					atomic.AddInt32(&lockFailures, 1)
				} else {
					t.Errorf("unexpected error for nonce %d: %v", nonce, err)
				}
			}
		}(int64(1000 + i))
	}
	wg.Wait()

	t.Logf("Number of lock failures: %d", atomic.LoadInt32(&lockFailures))

	var fromAcc orm.Account
	var toAcc orm.Account
	if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
		t.Fatalf("failed to query from account: %v", err)
	}
	if err := db.First(&toAcc, "user_id = ?", 2).Error; err != nil {
		t.Fatalf("failed to query to account: %v", err)
	}

	total := fromAcc.AvailableBalance.Add(fromAcc.ReservedBalance).Add(toAcc.AvailableBalance).Add(toAcc.ReservedBalance)
	expectedTotal := decimal.NewFromInt(1500)
	if !total.Equal(expectedTotal) {
		t.Errorf("total amount mismatch: expected %s, got %s", expectedTotal.String(), total.String())
	}

	if fromAcc.AvailableBalance.LessThan(decimal.Zero) {
		t.Errorf("overdraft occurred: available balance is negative: %s", fromAcc.AvailableBalance.String())
	}
}

func setupTestService(t *testing.T) (*gorm.DB, repository.TradeRepository, service.TradeService, *miniredis.Miniredis) {
	db := test.NewTestContainerDB(t)
	miniredis, redisClient := test.NewDummyRedis(t)
	test.SetupAccounts(t, db)

	repo := repository.NewTradeRepo()
	lockClient := distributedlock.NewRedisLockClient(redisClient)
	validator := NewTCCValidator(db, repo)
	service := NewTCCService(db, repo, validator, lockClient)

	return db, repo, service, miniredis
}

func setupMockTestService(t *testing.T, repo *mock.DummyTradeRepo) (*gorm.DB, repository.TradeRepository, service.TradeService, *miniredis.Miniredis) {
	db := test.NewTestContainerDB(t)
	miniredis, redisClient := test.NewDummyRedis(t)
	test.SetupAccounts(t, db)

	lockClient := distributedlock.NewRedisLockClient(redisClient)
	validator := NewTCCValidator(db, repo)
	service := NewTCCService(db, repo, validator, lockClient)

	return db, repo, service, miniredis
}
