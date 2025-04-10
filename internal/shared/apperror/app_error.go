package apperror

import (
	"errors"
	"fmt"
	"points/internal/shared/errcode"
)

type AppError struct {
	Code errcode.ErrorCode
	Msg  string
	Err  error
}

func (v *AppError) Error() string {
	if v.Err != nil {
		return fmt.Sprintf("%s: %s: %v", v.Code.String(), v.Msg, v.Err)
	}
	return fmt.Sprintf("%s: %s", v.Code.String(), v.Msg)
}

func (v *AppError) Unwrap() error {
	return v.Err
}

func Wrap(code errcode.ErrorCode, msg string, err error) error {
	return &AppError{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

func NewAppError(err error) *AppError {
	if err == nil {
		return &AppError{
			Code: errcode.ErrOK,
			Msg:  errcode.ErrOK.GetMessage(),
			Err:  nil,
		}
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return &AppError{
		Code: errcode.ErrInternal,
		Msg:  errcode.ErrInternal.GetMessage(),
		Err:  err,
	}
}
