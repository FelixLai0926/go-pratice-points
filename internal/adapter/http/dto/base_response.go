package dto

import "points/internal/shared/errcode"

type BaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewBaseResponse(status string, message string) *BaseResponse {
	return &BaseResponse{
		Status:  status,
		Message: message,
	}
}

func NewSuccessResponse() *BaseResponse {
	return &BaseResponse{
		Status:  errcode.ErrOK.String(),
		Message: errcode.ErrOK.GetMessage(),
	}
}
