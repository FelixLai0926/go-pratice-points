package errcode

import (
	"fmt"
	"net/http"
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
	ErrInsufficientBalance ErrorCode = 2002
	ErrUpdateAccount       ErrorCode = 2003
	ErrCreateTransaction   ErrorCode = 2004
	ErrGetTransaction      ErrorCode = 2005
	ErrUpdateTransaction   ErrorCode = 2006
	ErrPayloadMarshal      ErrorCode = 2007
	ErrCreateEvent         ErrorCode = 2008
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
	case ErrInsufficientBalance:
		return "insufficient balance"
	case ErrUpdateAccount:
		return "update account failed"
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
	default:
		return "unknown error"
	}
}

func (e ErrorCode) HTTPCode() int {
	switch e {
	case ErrOK:
		return http.StatusOK
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrInvalidRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrConflict:
		return http.StatusConflict
	case ErrGetAccount:
		return http.StatusBadRequest
	case ErrInsufficientBalance:
		return http.StatusBadRequest
	case ErrUpdateAccount:
		return http.StatusInternalServerError
	case ErrCreateTransaction:
		return http.StatusInternalServerError
	case ErrGetTransaction:
		return http.StatusBadRequest
	case ErrUpdateTransaction:
		return http.StatusInternalServerError
	case ErrPayloadMarshal:
		return http.StatusInternalServerError
	case ErrCreateEvent:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
