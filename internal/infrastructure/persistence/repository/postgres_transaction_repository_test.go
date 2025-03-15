package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/domain/valueobject"
	"points/internal/infrastructure"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/test"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewTradeRecordsRepo(db, config)

	ctx := context.Background()
	txRecord := &entity.TradeRecords{
		TransactionID: "tx1",
		Nonce:         1,
		FromAccountID: 100,
		ToAccountID:   200,
		Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
		Status:        1,
	}

	err := repoImpl.CreateTradeRecord(ctx, txRecord)
	assert.NoError(t, err, "error create trade record")

	var got model.TradeRecord
	err = db.Where("from_account_id = ? AND nonce = ?", 100, 1).First(&got).Error
	assert.NoError(t, err, "failed to query inserted transaction")
	assert.Equal(t, "tx1", got.TransactionID)
	assert.Equal(t, int32(1), got.Status)
	assert.Equal(t, txRecord.Amount.Value(), got.Amount, "amount should be equal")
}

func TestCreateOrUpdateTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewTradeRecordsRepo(db, config)
	ctx := context.Background()
	txRecord := &entity.TradeRecords{
		TransactionID: "tx1",
		Nonce:         1,
		FromAccountID: 100,
		ToAccountID:   200,
		Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
		Status:        1,
	}

	err := repoImpl.CreateOrUpdateTradeRecord(ctx, txRecord)
	assert.NoError(t, err, "error inserting transaction")

	var got model.TradeRecord
	err = db.Where("from_account_id = ? AND nonce = ?", 100, 1).First(&got).Error
	assert.NoError(t, err, "failed to query inserted transaction")
	assert.Equal(t, "tx1", got.TransactionID)
	assert.Equal(t, int32(1), got.Status)

	txRecord.Status = 2
	err = repoImpl.CreateOrUpdateTradeRecord(ctx, txRecord)
	assert.NoError(t, err, "error updating transaction")

	var updated model.TradeRecord
	err = db.Where("from_account_id = ? AND nonce = ?", 100, 1).First(&updated).Error
	assert.NoError(t, err, "failed to query updated transaction")
	assert.Equal(t, int32(2), updated.Status, "status should be updated to 2")
}

func TestGetTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewTradeRecordsRepo(db, config)
	ctx := context.Background()
	testNonce := int64(12345)
	testFrom := int64(1)
	testStatus := valueobject.TccPending

	trans := model.TradeRecord{
		Nonce:         testNonce,
		FromAccountID: testFrom,
		ToAccountID:   2,
		Amount:        decimal.NewFromInt(100),
		Status:        int32(testStatus),
	}
	if err := db.Create(&trans).Error; err != nil {
		t.Fatalf("failed to create test transaction: %v", err)
	}

	testCases := []struct {
		name           string
		nonce          int64
		from           int64
		status         *valueobject.TccStatus
		expectFound    bool
		expectedStatus int32
	}{
		{
			name:           "With status filter",
			nonce:          testNonce,
			from:           testFrom,
			status:         &testStatus,
			expectFound:    true,
			expectedStatus: int32(testStatus),
		},
		{
			name:           "Without status filter (nil)",
			nonce:          testNonce,
			from:           testFrom,
			status:         nil,
			expectFound:    true,
			expectedStatus: int32(testStatus),
		},
		{
			name:        "Transaction not found",
			nonce:       999,
			from:        testFrom,
			status:      &testStatus,
			expectFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotTrans, err := repoImpl.GetTradeRecord(ctx, tc.nonce, tc.from, tc.status)
			if !tc.expectFound {
				if err == nil {
					t.Errorf("expected error when transaction not found, but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetTransaction returned error: %v", err)
			}

			if gotTrans.FromAccountID != trans.FromAccountID || gotTrans.Nonce != trans.Nonce {
				t.Errorf("expected composite key (from_account_id, nonce): (%d, %d), got (%d, %d)",
					trans.FromAccountID, trans.Nonce, gotTrans.FromAccountID, gotTrans.Nonce)
			}
			if !gotTrans.Amount.Equals(valueobject.NewMoneyFromDecimal(trans.Amount)) {
				t.Errorf("expected Amount %s, got %s", trans.Amount.String(), gotTrans.Amount.String())
			}
			if gotTrans.Status != tc.expectedStatus {
				t.Errorf("expected Status %d, got %d", tc.expectedStatus, gotTrans.Status)
			}
		})
	}
}
