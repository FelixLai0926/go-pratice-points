package trade

import (
	"net/http"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/trade"

	"github.com/gin-gonic/gin"
)

func (h *TransferHandler) Confirm(c *gin.Context) {
	var req struct {
		From  *int32 `json:"from" form:"from" binding:"required"`
		To    *int32 `json:"to" form:"to" binding:"required"`
		Nonce *int64 `json:"nonce" form:"nonce" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transferRequest := &trade.ConfirmRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   c.Request.Context(),
			From:  *req.From,
			To:    *req.To,
			Nonce: *req.Nonce,
		},
	}

	if err := h.Validator.ValidateConfirmRequest(transferRequest); err != nil {
		c.Error(errors.NewAppError(http.StatusBadRequest, err))
		return
	}

	if err := h.Service.ManualConfirm(transferRequest); err != nil {
		c.Error(errors.NewAppError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errcode.ErrOK.String(),
		"message": "OK",
	})
}
