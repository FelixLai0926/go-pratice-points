package repository

import (
	"points/pkg/models/enum/tcc"
	"points/pkg/models/orm"
	"points/pkg/module/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()
	userId := int32(1)

	if err := repoImpl.CreateAccount(db, userId); err != nil {
		t.Fatalf("CreateAccount error: %v", err)
	}

	var gotAccount orm.Account
	if err := db.First(&gotAccount, "user_id = ?", userId).Error; err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
}

func TestGetAndUpdateAccount(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()

	account := orm.Account{
		UserID:           1,
		AvailableBalance: 100.00,
		ReservedBalance:  0,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	got, err := repoImpl.GetAccount(db, 1)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if got.AvailableBalance != 100.00 {
		t.Errorf("expected available balance 100, got %v", got.AvailableBalance)
	}

	got.AvailableBalance -= 20
	got.ReservedBalance += 20
	if err := repoImpl.UpdateAccount(db, got); err != nil {
		t.Fatalf("UpdateAccount error: %v", err)
	}

	updated, err := repoImpl.GetAccount(db, 1)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if updated.AvailableBalance != 80.00 || updated.ReservedBalance != 20 {
		t.Errorf("account not updated correctly: available = %v, reserved = %v", updated.AvailableBalance, updated.ReservedBalance)
	}
}

func TestCreateAndUpdateTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()

	trans := orm.Transaction{
		Nonce:         1,
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        50.00,
		Status:        0,
	}

	if err := repoImpl.CreateTransaction(db, &trans); err != nil {
		t.Fatalf("CreateTransaction error: %v", err)
	}

	trans.Status = 1
	if err := repoImpl.UpdateTransaction(db, &trans); err != nil {
		t.Fatalf("UpdateTransaction error: %v", err)
	}

	var gotTrans orm.Transaction
	err := db.Where("from_account_id = ? AND nonce = ?", trans.FromAccountID, trans.Nonce).First(&gotTrans).Error
	if err != nil {
		t.Fatalf("failed to query inserted/updated transaction: %v", err)
	}

	assert.Equal(t, trans.TransactionID, gotTrans.TransactionID)
	assert.Equal(t, int32(1), gotTrans.Status)
}

func TestCreateTransactionEvent(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()

	event := orm.TransactionEvent{
		TransactionID: "test-uuid",
		EventType:     "try",
		Payload:       `{"action":"try","amount":50}`,
	}
	if err := repoImpl.CreateTransactionEvent(db, &event); err != nil {
		t.Fatalf("CreateTransactionEvent error: %v", err)
	}

	var gotEvent orm.TransactionEvent
	if err := db.First(&gotEvent, "transaction_id = ?", event.TransactionID).Error; err != nil {
		t.Fatalf("failed to get transaction event: %v", err)
	}
	if gotEvent.EventType != "try" {
		t.Errorf("expected event type 'try', got %v", gotEvent.EventType)
	}
}

func TestCreateOrUpdateTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()

	txRecord := &orm.Transaction{
		TransactionID: "tx1",
		Nonce:         1,
		FromAccountID: 100,
		ToAccountID:   200,
		Amount:        100.0,
		Status:        1,
	}

	err := repoImpl.CreateOrUpdateTransaction(db, txRecord)
	assert.NoError(t, err, "error inserting transaction")

	var got orm.Transaction
	err = db.Where("from_account_id = ? AND nonce = ?", 100, 1).First(&got).Error
	assert.NoError(t, err, "failed to query inserted transaction")
	assert.Equal(t, "tx1", got.TransactionID)
	assert.Equal(t, int32(1), got.Status)

	txRecord.Status = 2
	err = repoImpl.CreateOrUpdateTransaction(db, txRecord)
	assert.NoError(t, err, "error updating transaction")

	var updated orm.Transaction
	err = db.Where("from_account_id = ? AND nonce = ?", 100, 1).First(&updated).Error
	assert.NoError(t, err, "failed to query updated transaction")
	assert.Equal(t, int32(2), updated.Status, "status should be updated to 2")
}

func TestGetTransaction(t *testing.T) {
	db := test.NewTestContainerDB(t)
	repoImpl := NewTradeRepo()

	testNonce := int64(12345)
	testFrom := int32(1)
	testStatus := tcc.Pending

	trans := orm.Transaction{
		Nonce:         testNonce,
		FromAccountID: testFrom,
		ToAccountID:   2,
		Amount:        100.0,
		Status:        int32(testStatus),
	}
	if err := db.Create(&trans).Error; err != nil {
		t.Fatalf("failed to create test transaction: %v", err)
	}

	testCases := []struct {
		name           string
		nonce          int64
		from           int32
		status         *tcc.Status
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
			gotTrans, err := repoImpl.GetTransaction(db, tc.nonce, tc.from, tc.status)
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
			if gotTrans.Amount != trans.Amount {
				t.Errorf("expected Amount %f, got %f", trans.Amount, gotTrans.Amount)
			}
			if gotTrans.Status != tc.expectedStatus {
				t.Errorf("expected Status %d, got %d", tc.expectedStatus, gotTrans.Status)
			}
		})
	}
}
