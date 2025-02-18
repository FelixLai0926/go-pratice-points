package eventpayload

import "encoding/json"

type TransferPayload struct {
	Action string  `json:"action"`
	Amount float64 `json:"amount"`
}

func (p TransferPayload) ToJSON() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
