package repository

import (
	"context"
	"points/internal/domain/valueobject"
	"points/internal/infrastructure"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/test"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCreateAccount(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewAccountRepo(db, config)
	userId := int64(1)
	ctx := context.Background()

	if err := repoImpl.CreateAccount(ctx, userId); err != nil {
		t.Fatalf("CreateAccount error: %v", err)
	}

	var gotAccount model.Account
	if err := db.First(&gotAccount, "user_id = ?", userId).Error; err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
}

func TestReserveBalance(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewAccountRepo(db, config)
	ctx := context.Background()
	userId := int64(1)

	account := model.Account{
		UserID:           userId,
		AvailableBalance: decimal.NewFromInt(100),
		ReservedBalance:  decimal.Zero,
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	if err := repoImpl.ReserveBalance(ctx, userId, valueobject.NewMoneyFromDecimal(decimal.NewFromInt(20))); err != nil {
		t.Fatalf("UpdateAccount error: %v", err)
	}

	updated, err := repoImpl.GetAccount(ctx, userId)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if !updated.AvailableBalance.Equals(valueobject.NewMoneyFromDecimal(decimal.NewFromInt(80))) || !updated.ReservedBalance.Equals(valueobject.NewMoneyFromDecimal(decimal.NewFromInt(20))) {
		t.Errorf("account not updated correctly: available = %v, reserved = %v", updated.AvailableBalance, updated.ReservedBalance)
	}
}

func TestUnreserveBalance(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewAccountRepo(db, config)
	ctx := context.Background()
	from := int64(1)
	to := int64(2)

	accountFrom := model.Account{
		UserID:           from,
		AvailableBalance: decimal.Zero,
		ReservedBalance:  decimal.NewFromInt(100.00),
	}
	if err := db.Create(&accountFrom).Error; err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accountTo := model.Account{
		UserID:           to,
		AvailableBalance: decimal.Zero,
		ReservedBalance:  decimal.Zero,
	}
	if err := db.Create(&accountTo).Error; err != nil {
		t.Fatalf("failed to create account to: %v", err)
	}

	if err := repoImpl.UnreserveBalance(ctx, from, to, valueobject.NewMoneyFromDecimal(decimal.NewFromInt(20))); err != nil {
		t.Fatalf("UpdateAccount error: %v", err)
	}

	fromAccount, err := repoImpl.GetAccount(ctx, from)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if !fromAccount.ReservedBalance.Equals(valueobject.NewMoneyFromDecimal(decimal.NewFromInt(80))) || !fromAccount.AvailableBalance.Equals(valueobject.Zero) {
		t.Errorf("account not updated correctly: available = %v, reserved = %v", fromAccount.AvailableBalance, fromAccount.ReservedBalance)
	}
	toAccount, err := repoImpl.GetAccount(ctx, to)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}
	if !toAccount.AvailableBalance.Equals(valueobject.NewMoneyFromDecimal(decimal.NewFromInt(20))) {
		t.Errorf("account not updated correctly: available = %v, reserved = %v", fromAccount.AvailableBalance, fromAccount.ReservedBalance)
	}
}

func TestGetAccount(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewAccountRepo(db, config)
	ctx := context.Background()
	userId := int64(42)

	account := model.Account{
		UserID:           userId,
		AvailableBalance: decimal.NewFromInt(150),
		ReservedBalance:  decimal.NewFromInt(50),
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	gotAccount, err := repoImpl.GetAccount(ctx, userId)
	if err != nil {
		t.Fatalf("GetAccount error: %v", err)
	}

	expectedAvail := valueobject.NewMoneyFromDecimal(decimal.NewFromInt(150))
	expectedReserved := valueobject.NewMoneyFromDecimal(decimal.NewFromInt(50))

	if gotAccount.UserID != userId {
		t.Errorf("expected user id %d, got %d", userId, gotAccount.UserID)
	}
	if !gotAccount.AvailableBalance.Equals(expectedAvail) {
		t.Errorf("expected available balance %v, got %v", expectedAvail, gotAccount.AvailableBalance)
	}
	if !gotAccount.ReservedBalance.Equals(expectedReserved) {
		t.Errorf("expected reserved balance %v, got %v", expectedReserved, gotAccount.ReservedBalance)
	}
}
