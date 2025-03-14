package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"points/internal/domain/entity"
	"points/internal/domain/valueobject"
	"points/test/mock"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func dummyAccount(userID int64, availableBalance, reservedBalance decimal.Decimal) *entity.Account {
	return &entity.Account{
		UserID:           userID,
		AvailableBalance: valueobject.NewMoneyFromDecimal(availableBalance),
		ReservedBalance:  valueobject.NewMoneyFromDecimal(reservedBalance),
	}
}

func TestTransferTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		setupMocks func(
			uow *mock.MockUnitOfWork,
			accRepo *mock.MockAccountRepository,
			transRepo *mock.MockTradeRecordsRepository,
			eventRepo *mock.MockTransactionEventRepository,
		)
		expectedErr error
	}{
		{
			name: "success - account not exists, create account, all ok",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(nil, nil).Times(1)
				accRepo.EXPECT().CreateAccount(ctx, int64(2)).Return(nil).Times(1)
				accRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), nil).Return(nil, nil).Times(1)
				accRepo.EXPECT().ReserveBalance(ctx, int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().CreateTradeRecord(ctx, gomock.AssignableToTypeOf(&entity.TradeRecords{})).
					DoAndReturn(func(ctx context.Context, tr *entity.TradeRecords) error {
						if tr.Status != int32(valueobject.TccPending) {
							return errors.New("invalid status")
						}
						return nil
					}).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.AssignableToTypeOf(&entity.TransactionEvent{})).
					Return(nil).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: nil,
		},
		{
			name: "fail - GetAccount error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(nil, errors.New("db error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("db error"),
		},
		{
			name: "fail - CreateAccount error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(nil, nil).Times(1)
				accRepo.EXPECT().CreateAccount(ctx, int64(2)).Return(errors.New("create error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("create error"),
		},
		{
			name: "fail - ReserveBalance error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				accRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), nil).Return(nil, nil).Times(1)
				accRepo.EXPECT().ReserveBalance(ctx, int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).
					Return(errors.New("reserve error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("reserve error"),
		},
		{
			name: "fail - CreateTradeRecord error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), nil).Return(nil, nil).Times(1)
				accRepo.EXPECT().ReserveBalance(ctx, int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).
					Return(nil).Times(1)
				transRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).Return(errors.New("create transaction error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
			},
			expectedErr: errors.New("create transaction error"),
		},
		{
			name: "fail - CreateTransactionEvent error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(1)).Return(dummyAccount(1, decimal.NewFromInt(100), decimal.Zero), nil).Times(1)
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), nil).Return(nil, nil).Times(1)
				accRepo.EXPECT().ReserveBalance(ctx, int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).
					Return(nil).Times(1)
				transRepo.EXPECT().CreateTradeRecord(ctx, gomock.Any()).Return(nil).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(errors.New("create event error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: errors.New("create event error"),
		},
		{
			name: "fail - conflict nonce error",
			setupMocks: func(uow *mock.MockUnitOfWork,
				accRepo *mock.MockAccountRepository,
				transRepo *mock.MockTradeRecordsRepository,
				eventRepo *mock.MockTransactionEventRepository) {
				accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(dummyAccount(2, decimal.Zero, decimal.Zero), nil).Times(1)
				conflictRecord := &entity.TradeRecords{TransactionID: "existing"}
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), nil).
					Return(conflictRecord, errors.New("conflict error")).Times(1)
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
			},
			expectedErr: errors.New("conflict nonce"),
		},
	}

	svc := NewTransactionApplicationService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uow := mock.NewMockUnitOfWork(ctrl)
			accRepo := mock.NewMockAccountRepository(ctrl)
			transRepo := mock.NewMockTradeRecordsRepository(ctrl)
			eventRepo := mock.NewMockTransactionEventRepository(ctrl)

			tt.setupMocks(uow, accRepo, transRepo, eventRepo)

			err := svc.TransferTransaction(ctx, uow, 123, 1, 2, valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)))
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfirmTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		setupMocks  func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository)
		expectedErr error
	}{
		{
			name: "success - confirm transaction",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}
				trans.Confirm()

				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(2), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(nil).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(2)

				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: nil,
		},
		{
			name: "fail - GetTradeRecord error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(nil, errors.New("get transaction error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
			},
			expectedErr: errors.New("get transaction error"),
		},
		{
			name: "fail - UnreserveBalance error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Confirm()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(2), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(errors.New("unreserve error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("unreserve error"),
		},
		{
			name: "fail - UpdateTradeRecord error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Confirm()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(2), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(errors.New("update error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("update error"),
		},
		{
			name: "fail - CreateTransactionEvent error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Confirm()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(2), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(nil).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(errors.New("create event error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: errors.New("create event error"),
		},
	}

	svc := NewTransactionApplicationService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uow := mock.NewMockUnitOfWork(ctrl)
			accRepo := mock.NewMockAccountRepository(ctrl)
			transRepo := mock.NewMockTradeRecordsRepository(ctrl)
			eventRepo := mock.NewMockTransactionEventRepository(ctrl)

			tt.setupMocks(uow, accRepo, transRepo, eventRepo)
			err := svc.ConfirmTransaction(ctx, uow, 123, 1, 2)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCancelTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		setupMocks  func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository)
		expectedErr error
	}{
		{
			name: "success - cancel transaction",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Cancel()

				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(nil).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(nil).Times(2)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: nil,
		},
		{
			name: "fail - GetTradeRecord error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(nil, errors.New("get transaction error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
			},
			expectedErr: errors.New("get transaction error"),
		},
		{
			name: "fail - UnreserveBalance error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Cancel()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(errors.New("unreserve error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("unreserve error"),
		},
		{
			name: "fail - UpdateTradeRecord error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Cancel()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(errors.New("update error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
			},
			expectedErr: errors.New("update error"),
		},
		{
			name: "fail - CreateTransactionEvent error",
			setupMocks: func(uow *mock.MockUnitOfWork, accRepo *mock.MockAccountRepository, transRepo *mock.MockTradeRecordsRepository, eventRepo *mock.MockTransactionEventRepository) {
				pendingStatus := valueobject.TccPending
				trans := &entity.TradeRecords{
					TransactionID: "tx-123",
					FromAccountID: 1,
					ToAccountID:   2,
					Amount:        valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100)),
					Status:        int32(pendingStatus),
				}
				trans.Cancel()
				transRepo.EXPECT().GetTradeRecord(ctx, int64(123), int64(1), &pendingStatus).Return(trans, nil).Times(1)
				accRepo.EXPECT().UnreserveBalance(ctx, int64(1), int64(1), valueobject.NewMoneyFromDecimal(decimal.NewFromInt(100))).Return(nil).Times(1)
				transRepo.EXPECT().UpdateTradeRecord(ctx, trans).Return(nil).Times(1)
				eventRepo.EXPECT().CreateTransactionEvent(ctx, gomock.Any()).Return(errors.New("create event error")).Times(1)
				uow.EXPECT().TradeRecordsRepository().Return(transRepo).AnyTimes()
				uow.EXPECT().AccountRepository().Return(accRepo).AnyTimes()
				uow.EXPECT().TransactionEventRepository().Return(eventRepo).AnyTimes()
			},
			expectedErr: errors.New("create event error"),
		},
	}

	svc := NewTransactionApplicationService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uow := mock.NewMockUnitOfWork(ctrl)
			accRepo := mock.NewMockAccountRepository(ctrl)
			transRepo := mock.NewMockTradeRecordsRepository(ctrl)
			eventRepo := mock.NewMockTransactionEventRepository(ctrl)

			tt.setupMocks(uow, accRepo, transRepo, eventRepo)
			err := svc.CancelTransaction(ctx, uow, 123, 1, 2)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
