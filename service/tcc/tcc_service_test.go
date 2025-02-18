package tcc

import (
	"context"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"
	"points/pkg/models/trade"
	"points/pkg/module/distributedlock"
	"points/pkg/module/test"
	"points/repository"
	"points/service"
	"strings"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"gorm.io/gorm"
)

func TestTransfer_AutoConfirm(t *testing.T) {
	db, _, service, miniredis := setupTestService(t)
	defer miniredis.Close()

	req := &trade.TransferRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   context.Background(),
			From:  1,
			To:    2,
			Nonce: 12345,
		},
		Amount:      100.0,
		AutoConfirm: true,
	}

	if err := service.Transfer(req); err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}

	var fromAcc orm.Account
	if err := db.First(&fromAcc, "user_id = ?", req.From).Error; err != nil {
		t.Fatalf("failed to query from account: %v", err)
	}
	if fromAcc.AvailableBalance != 900.0 {
		t.Errorf("expected from account available balance 900, got %f", fromAcc.AvailableBalance)
	}
	if fromAcc.ReservedBalance != 0.0 {
		t.Errorf("expected from account reserved balance 0, got %f", fromAcc.ReservedBalance)
	}

	var toAcc orm.Account
	if err := db.First(&toAcc, "user_id = ?", req.To).Error; err != nil {
		t.Fatalf("failed to query to account: %v", err)
	}
	if toAcc.AvailableBalance != 600.0 {
		t.Errorf("expected to account available balance 600, got %f", toAcc.AvailableBalance)
	}
}

func TestCancel_PendingTransfer(t *testing.T) {
	db, repo, service, miniredis := setupTestService(t)
	defer miniredis.Close()

	transferReq := &trade.TransferRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   context.Background(),
			From:  1,
			To:    2,
			Nonce: 54321,
		},
		Amount:      50.0,
		AutoConfirm: false,
	}

	if err := service.Transfer(transferReq); err != nil {
		t.Fatalf("Transfer (Try phase) failed: %v", err)
	}

	var fromAcc orm.Account
	if err := db.First(&fromAcc, "user_id = ?", transferReq.From).Error; err != nil {
		t.Fatalf("failed to query from account: %v", err)
	}
	if fromAcc.ReservedBalance != 50.0 {
		t.Errorf("expected reserved balance 50, got %f", fromAcc.ReservedBalance)
	}

	cancelReq := &trade.CancelRequest{
		BaseRequest: transferReq.BaseRequest,
	}

	if err := service.Cancel(cancelReq); err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}

	if err := db.First(&fromAcc, "user_id = ?", cancelReq.From).Error; err != nil {
		t.Fatalf("failed to query from account after cancel: %v", err)
	}
	if fromAcc.ReservedBalance != 0.0 {
		t.Errorf("after cancel, expected reserved balance 0, got %f", fromAcc.ReservedBalance)
	}
	if fromAcc.AvailableBalance != 1000.0 {
		t.Errorf("after cancel, expected available balance 1000, got %f", fromAcc.AvailableBalance)
	}

	pendingStatus := tcc.Pending
	trans, err := repo.GetTransaction(db, cancelReq.Nonce, cancelReq.From, &pendingStatus)
	if err == nil {
		if trans.Status != int32(tcc.Canceled) {
			t.Errorf("expected transaction status Canceled, got %d", trans.Status)
		}
	}
}

func TestConcurrentTransfers(t *testing.T) {
	db, _, service, miniredis := setupTestService(t)
	defer miniredis.Close()

	numTransfers := 10
	transferAmount := 100.0

	var wg sync.WaitGroup
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
			if err := service.Transfer(req); err != nil {
				t.Errorf("transfer failed for nonce %d: %v", nonce, err)
			}
		}(int64(i + 1000))
	}
	wg.Wait()

	var fromAcc orm.Account
	if err := db.First(&fromAcc, "user_id = ?", 1).Error; err != nil {
		t.Fatalf("failed to query from account: %v", err)
	}

	expectedBalance := 1000.0 - float64(numTransfers)*transferAmount
	if fromAcc.AvailableBalance != expectedBalance {
		t.Errorf("expected from account available balance %f, got %f", expectedBalance, fromAcc.AvailableBalance)
	}

	if fromAcc.AvailableBalance < 0 {
		t.Errorf("overdraft occurred: available balance is negative: %f", fromAcc.AvailableBalance)
	}

	var toAcc orm.Account
	if err := db.First(&toAcc, "user_id = ?", 2).Error; err != nil {
		t.Fatalf("failed to query to account: %v", err)
	}

	expectedToBalance := 500.0 + float64(numTransfers)*transferAmount
	if toAcc.AvailableBalance != expectedToBalance {
		t.Errorf("expected to account available balance %f, got %f", expectedToBalance, toAcc.AvailableBalance)
	}
}

func TestEnsureDestinationAccount(t *testing.T) {
	db, _, service, _ := setupTestService(t)

	tests := []struct {
		name        string
		userID      int32
		setup       func(t *testing.T, db *gorm.DB)
		wantExists  bool
		wantErr     bool
		expectedErr string
	}{
		{
			name:        "Account already exists",
			userID:      1,
			setup:       func(t *testing.T, db *gorm.DB) {},
			wantExists:  true,
			wantErr:     false,
			expectedErr: "",
		},
		{
			name:   "Account does not exist, creation succeeds",
			userID: 3,
			setup: func(t *testing.T, db *gorm.DB) {
				if err := db.Exec("DELETE FROM account WHERE user_id = ?", 3).Error; err != nil {
					t.Fatalf("failed to delete account: %v", err)
				}
			},
			wantExists:  true,
			wantErr:     false,
			expectedErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t, db)
			}

			rq := &trade.EnsureAccountRequest{
				Ctx:    context.Background(),
				UserID: tc.userID,
			}
			err := service.EnsureDestinationAccount(rq)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("expected error to contain %q, got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}

			var account orm.Account
			err = db.First(&account, "user_id = ?", tc.userID).Error
			if tc.wantExists {
				if err != nil {
					t.Fatalf("expected account to exist for user_id %d, but not found: %v", tc.userID, err)
				}
			} else {
				if err == nil {
					t.Errorf("expected account not to exist for user_id %d, but found one", tc.userID)
				}
			}
		})
	}
}

func setupTestService(t *testing.T) (*gorm.DB, repository.TradeRepository, service.TradeService, *miniredis.Miniredis) {
	db := test.NewTestContainerDB(t)
	miniredis, redisClient := test.NewDummyRedis(t)
	test.SetupAccounts(t, db)

	repo := repository.NewTradeRepo()
	lockClient := distributedlock.NewRedisLockClient(redisClient)
	service := NewTCCService(db, repo, lockClient)

	return db, repo, service, miniredis
}
