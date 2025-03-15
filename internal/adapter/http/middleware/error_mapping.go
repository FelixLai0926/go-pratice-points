package middleware

import (
	"net/http"
	"points/internal/shared/errcode"
)

func mapErrorCodeToHTTPStatus(e errcode.ErrorCode) int {
	switch e {
	case errcode.ErrOK:
		return http.StatusOK
	case errcode.ErrInternal:
		return http.StatusInternalServerError
	case errcode.ErrInvalidRequest:
		return http.StatusBadRequest
	case errcode.ErrNotFound:
		return http.StatusBadRequest
	case errcode.ErrUnauthorized:
		return http.StatusUnauthorized
	case errcode.ErrConflict:
		return http.StatusConflict
	case errcode.ErrGetAccount:
		return http.StatusBadRequest
	case errcode.ErrCreateAccount:
		return http.StatusInternalServerError
	case errcode.ErrAccountNotFound:
		return http.StatusNotFound
	case errcode.ErrInsufficientBalance:
		return http.StatusBadRequest
	case errcode.ErrReserveBalance:
		return http.StatusInternalServerError
	case errcode.ErrUnreserveBalance:
		return http.StatusInternalServerError
	case errcode.ErrCreateTransaction:
		return http.StatusInternalServerError
	case errcode.ErrGetTransaction:
		return http.StatusBadRequest
	case errcode.ErrUpdateTransaction:
		return http.StatusInternalServerError
	case errcode.ErrPayloadMarshal:
		return http.StatusInternalServerError
	case errcode.ErrCreateEvent:
		return http.StatusInternalServerError
	case errcode.ErrDistrubutedLockNotObtained:
		return http.StatusInternalServerError
	case errcode.ErrDistrubutedLockAcquire:
		return http.StatusInternalServerError
	case errcode.ErrDistrubutedLockRelease:
		return http.StatusInternalServerError
	case errcode.ErrDistrubutedLockRenew:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
