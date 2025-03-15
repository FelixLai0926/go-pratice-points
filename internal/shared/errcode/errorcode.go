package errcode

import (
	"fmt"
)

type ErrorCode int

const (
	ErrOK             ErrorCode = 0
	ErrInternal       ErrorCode = 1000
	ErrInvalidRequest ErrorCode = 1001
	ErrNotFound       ErrorCode = 1002
	ErrUnauthorized   ErrorCode = 1003
	ErrConflict       ErrorCode = 1004

	ErrGetAccount          ErrorCode = 2001
	ErrCreateAccount       ErrorCode = 2002
	ErrAccountNotFound     ErrorCode = 2003
	ErrInsufficientBalance ErrorCode = 2004
	ErrReserveBalance      ErrorCode = 2005
	ErrUnreserveBalance    ErrorCode = 2006
	ErrCreateTransaction   ErrorCode = 2007
	ErrGetTransaction      ErrorCode = 2008
	ErrUpdateTransaction   ErrorCode = 2009
	ErrPayloadMarshal      ErrorCode = 2010
	ErrCreateEvent         ErrorCode = 2011

	ErrDistrubutedLockNotObtained ErrorCode = 3001
	ErrDistrubutedLockAcquire     ErrorCode = 3002
	ErrDistrubutedLockRelease     ErrorCode = 3003
	ErrDistrubutedLockRenew       ErrorCode = 3004
)

func (e ErrorCode) String() string {
	return fmt.Sprintf("%04d", int(e))
}

func (e ErrorCode) GetMessage() string {
	switch e {
	case ErrOK:
		return "OK"
	case ErrInternal:
		return "internal error"
	case ErrInvalidRequest:
		return "invalid request"
	case ErrNotFound:
		return "not found"
	case ErrUnauthorized:
		return "unauthorized"
	case ErrConflict:
		return "conflict"
	case ErrGetAccount:
		return "get account failed"
	case ErrCreateAccount:
		return "create account failed"
	case ErrAccountNotFound:
		return "account not found"
	case ErrInsufficientBalance:
		return "insufficient balance"
	case ErrReserveBalance:
		return "reserve balance failed"
	case ErrUnreserveBalance:
		return "unreserve balance failed"
	case ErrCreateTransaction:
		return "create transaction failed"
	case ErrGetTransaction:
		return "get transaction failed"
	case ErrUpdateTransaction:
		return "update transaction failed"
	case ErrPayloadMarshal:
		return "payload marshal failed"
	case ErrCreateEvent:
		return "create event failed"
	case ErrDistrubutedLockNotObtained:
		return "distributed lock not obtained"
	case ErrDistrubutedLockAcquire:
		return "distributed lock acquire failed"
	case ErrDistrubutedLockRelease:
		return "distributed lock release failed"
	case ErrDistrubutedLockRenew:
		return "distributed lock renew failed"
	default:
		return "unknown error"
	}
}
