package tcc

import (
	"context"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"
	"points/pkg/models/trade"
	"points/pkg/module/test"
	"points/repository"
	"points/service"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func TestValidateRequest(t *testing.T) {
	db, validator := setupTestValidator(t)

	existingTrans := orm.Transaction{
		Nonce:         123,
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        decimal.NewFromInt(100),
		Status:        int32(tcc.Pending),
		TransactionID: "tx-123",
	}
	if err := db.Create(&existingTrans).Error; err != nil {
		t.Fatalf("failed to create existing transaction: %v", err)
	}

	testCases := []struct {
		name          string
		request       *trade.TransferRequest
		expectErrText string
	}{
		{
			name: "Account not found",
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  999,
					To:    2,
					Nonce: 1,
				},
				Amount:      decimal.NewFromInt(100.0),
				AutoConfirm: true,
			},
			expectErrText: "source account not found",
		},
		{
			name: "Nonce already used",
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 123,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: true,
			},
			expectErrText: "nonce already used",
		},
		{
			name: "Valid request",
			request: &trade.TransferRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					To:    2,
					Nonce: 456,
				},
				Amount:      decimal.NewFromInt(100),
				AutoConfirm: true,
			},
			expectErrText: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateTransferRequest(tc.request)
			if tc.expectErrText != "" {
				if err == nil {
					t.Fatalf("expected error containing %q but got nil", tc.expectErrText)
				}
				if !strings.Contains(err.Error(), tc.expectErrText) {
					t.Fatalf("expected error containing %q but got: %v", tc.expectErrText, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateConfirmRequest(t *testing.T) {
	db, validator := setupTestValidator(t)

	tests := []struct {
		name          string
		req           *trade.ConfirmRequest
		setup         func(db *gorm.DB)
		wantErr       bool
		expectedError string
		wantErrCode   errcode.ErrorCode
	}{
		{
			name: "Account not found",
			req: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  999,
					Nonce: 1,
				},
			},
			setup:         func(db *gorm.DB) {},
			wantErr:       true,
			expectedError: "source account not found",
			wantErrCode:   errcode.ErrNotFound,
		},
		{
			name: "Transaction not found",
			req: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 12345,
				},
			},
			setup:         func(db *gorm.DB) {},
			wantErr:       true,
			expectedError: "transaction not found",
			wantErrCode:   errcode.ErrGetTransaction,
		},
		{
			name: "Insufficient reserved balance",
			req: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 789,
				},
			},
			setup: func(db *gorm.DB) {
				var account orm.Account
				if err := db.First(&account, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query account: %v", err)
				}
				account.ReservedBalance.Equal(decimal.NewFromInt(50))
				if err := db.Save(&account).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				pendingTx := orm.Transaction{
					Nonce:         789,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(100),
					Status:        int32(tcc.Pending),
					TransactionID: "tx-789",
				}
				if err := db.Create(&pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
			},
			wantErr:       true,
			expectedError: "insufficient balance",
			wantErrCode:   errcode.ErrInsufficientBalance,
		},
		{
			name: "Valid confirm request",
			req: &trade.ConfirmRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 456,
					To:    2,
				},
			},
			setup: func(db *gorm.DB) {
				var account orm.Account
				if err := db.First(&account, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query account: %v", err)
				}
				account.ReservedBalance = decimal.NewFromInt(150)
				if err := db.Save(&account).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				pendingTx := orm.Transaction{
					Nonce:         456,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(100),
					Status:        int32(tcc.Pending),
					TransactionID: "tx-456",
				}
				if err := db.Create(&pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(db)
			}
			err := validator.ValidateConfirmRequest(tc.req)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q but got nil", tc.expectedError)
				}
				appErr := err.(*errors.AppError)
				if appErr.Code != tc.wantErrCode {
					t.Errorf("expected error code %s, got %s", tc.wantErrCode.String(), appErr.Code.String())
				}

				if !strings.Contains(appErr.Error(), tc.expectedError) {
					t.Errorf("expected error message to contain %q, got %q", tc.expectedError, appErr.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateCancelRequest(t *testing.T) {
	db, validator := setupTestValidator(t)

	tests := []struct {
		name          string
		req           *trade.CancelRequest
		setup         func(db *gorm.DB)
		wantErr       bool
		expectedError string
		wantErrCode   errcode.ErrorCode
	}{
		{
			name: "Account not found",
			req: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  999,
					Nonce: 1,
				},
			},
			setup:         func(db *gorm.DB) {},
			wantErr:       true,
			expectedError: "source account not found",
			wantErrCode:   errcode.ErrNotFound,
		},
		{
			name: "Transaction not found",
			req: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 12345,
				},
			},
			setup:         func(db *gorm.DB) {},
			wantErr:       true,
			expectedError: "transaction not found",
			wantErrCode:   errcode.ErrGetTransaction,
		},
		{
			name: "Insufficient reserved balance",
			req: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 789,
				},
			},
			setup: func(db *gorm.DB) {
				var account orm.Account
				if err := db.First(&account, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query account: %v", err)
				}
				account.ReservedBalance = decimal.NewFromInt(50)
				if err := db.Save(&account).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				pendingTx := orm.Transaction{
					Nonce:         789,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(100),
					Status:        int32(tcc.Pending),
					TransactionID: "tx-789",
				}
				if err := db.Create(&pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
			},
			wantErr:       true,
			expectedError: "insufficient balance",
			wantErrCode:   errcode.ErrInsufficientBalance,
		},
		{
			name: "Valid confirm request",
			req: &trade.CancelRequest{
				BaseRequest: trade.BaseRequest{
					Ctx:   context.Background(),
					From:  1,
					Nonce: 456,
					To:    2,
				},
			},
			setup: func(db *gorm.DB) {
				var account orm.Account
				if err := db.First(&account, "user_id = ?", 1).Error; err != nil {
					t.Fatalf("failed to query account: %v", err)
				}
				account.ReservedBalance = decimal.NewFromInt(150)
				if err := db.Save(&account).Error; err != nil {
					t.Fatalf("failed to update account: %v", err)
				}
				pendingTx := orm.Transaction{
					Nonce:         456,
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        decimal.NewFromInt(100),
					Status:        int32(tcc.Pending),
					TransactionID: "tx-456",
				}
				if err := db.Create(&pendingTx).Error; err != nil {
					t.Fatalf("failed to create pending transaction: %v", err)
				}
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(db)
			}
			err := validator.ValidateCancelRequest(tc.req)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q but got nil", tc.expectedError)
				}
				appErr := err.(*errors.AppError)
				if appErr.Code != tc.wantErrCode {
					t.Errorf("expected error code %s, got %s", tc.wantErrCode.String(), appErr.Code.String())
				}

				if !strings.Contains(appErr.Error(), tc.expectedError) {
					t.Errorf("expected error message to contain %q, got %q", tc.expectedError, appErr.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func setupTestValidator(t *testing.T) (*gorm.DB, service.TradeValidator) {
	db := test.NewTestContainerDB(t)
	test.SetupAccounts(t, db)

	repo := repository.NewTradeRepo()
	validator := NewTCCValidator(db, repo)

	return db, validator
}
