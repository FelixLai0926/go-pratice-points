package trade

import (
	"points/service"

	"gorm.io/gorm"
)

type TransferHandler struct {
	DB      *gorm.DB
	Service service.TradeService
}

func NewTransferHandler(db *gorm.DB, service service.TradeService) *TransferHandler {
	return &TransferHandler{
		DB:      db,
		Service: service,
	}
}
