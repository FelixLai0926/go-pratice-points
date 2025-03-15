package repository

import (
	"context"
	"points/internal/domain/entity"
	"points/internal/infrastructure"
	"points/internal/infrastructure/persistence/gorm/model"
	"points/test"
	"testing"
)

func TestCreateTransactionEvent(t *testing.T) {
	db := test.NewTestContainerDB(t)
	copier := infrastructure.NewCopierImpl()
	config := infrastructure.NewConfigImpl(nil, nil, copier)
	repoImpl := NewTransactionEventRepo(db, config)
	ctx := context.Background()
	event := entity.TransactionEvent{
		TransactionID: "test-uuid",
		EventType:     "try",
		Payload:       `{"action":"try","amount":50}`,
	}
	if err := repoImpl.CreateTransactionEvent(ctx, &event); err != nil {
		t.Fatalf("CreateTransactionEvent error: %v", err)
	}

	var gotEvent model.TransactionEvent
	if err := db.First(&gotEvent, "transaction_id = ?", event.TransactionID).Error; err != nil {
		t.Fatalf("failed to get transaction event: %v", err)
	}
	if gotEvent.EventType != "try" {
		t.Errorf("expected event type 'try', got %v", gotEvent.EventType)
	}
}
