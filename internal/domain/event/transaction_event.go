package event

import "points/internal/domain/valueobject"

const TransactionEventType = "TransactionEvent"

type TransactionEvent struct {
	TransactionID string
	Action        string
	FromAccountID int64
	ToAccountID   int64
	Amount        valueobject.Money
}

func (e TransactionEvent) EventType() string {
	return TransactionEventType
}
