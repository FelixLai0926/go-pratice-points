package tcc

type Status int32

const (
	Pending Status = iota
	Confirmed
	Canceled
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "pending"
	case Confirmed:
		return "confirmed"
	case Canceled:
		return "canceled"
	default:
		return "unknown"
	}
}
