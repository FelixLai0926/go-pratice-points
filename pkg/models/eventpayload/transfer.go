package eventpayload

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type TransferPayload struct {
	Action string          `json:"action"`
	Amount decimal.Decimal `json:"amount"`
}

func (p TransferPayload) ToJSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
