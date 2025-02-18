package errors

import (
	"fmt"
	"points/pkg/models/enum/errcode"
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
