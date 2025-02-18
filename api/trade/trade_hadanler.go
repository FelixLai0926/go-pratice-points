package trade

import (
	"points/service"

	"gorm.io/gorm"
)

type TransferHandler struct {
	DB        *gorm.DB
	Service   service.TradeService
	Validator service.TradeValidator
}

func NewTransferHandler(db *gorm.DB, service service.TradeService, validator service.TradeValidator) *TransferHandler {
	return &TransferHandler{
		DB:        db,
		Service:   service,
		Validator: validator,
	}
}
