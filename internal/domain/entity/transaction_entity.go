package entity

import (
	"points/internal/domain/event"
	"points/internal/domain/valueobject"
	"time"
)

type TradeRecords struct {
	TransactionID string
	Nonce         int64
	FromAccountID int64
	ToAccountID   int64
	Amount        valueobject.Money
	Status        int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	events        []event.TransactionEvent
}

func (t *TradeRecords) Transfer() {
	t.Status = int32(valueobject.TccPending)
	t.events = append(t.events, event.TransactionEvent{
		TransactionID: t.TransactionID,
		Action:        valueobject.TccPending.String(),
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
	})
}

func (t *TradeRecords) Confirm() {
	t.Status = int32(valueobject.TccConfirmed)
	t.events = append(t.events, event.TransactionEvent{
		TransactionID: t.TransactionID,
		Action:        valueobject.TccConfirmed.String(),
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
	})
}

func (t *TradeRecords) Cancel() {
	t.Status = int32(valueobject.TccCanceled)
	t.events = append(t.events, event.TransactionEvent{
		TransactionID: t.TransactionID,
		Action:        valueobject.TccCanceled.String(),
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
	})
}

func (t *TradeRecords) PullEvents() []event.TransactionEvent {
	evts := t.events
	t.events = []event.TransactionEvent{}
	return evts
}
