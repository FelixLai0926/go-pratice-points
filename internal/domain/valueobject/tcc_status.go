package valueobject

type TccStatus int32

const (
	TccPending TccStatus = iota
	TccConfirmed
	TccCanceled
)

func (s TccStatus) String() string {
	switch s {
	case TccPending:
		return "pending"
	case TccConfirmed:
		return "confirmed"
	case TccCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

func (s TccStatus) Ptr() *TccStatus {
	return &s
}
