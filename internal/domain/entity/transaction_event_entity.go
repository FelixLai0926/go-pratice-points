package entity

import (
	"time"
)

type TransactionEvent struct {
	ID            int32
	TransactionID string
	EventType     string
	Payload       string
	CreatedAt     time.Time
}
