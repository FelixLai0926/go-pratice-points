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
		From  *int64 `json:"from" form:"from" binding:"required"`
		To    *int64 `json:"to" form:"to" binding:"required"`
		Nonce *int64 `json:"nonce" form:"nonce" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(errors.Wrap(errcode.ErrInvalidRequest, "invalid request", err))
		return
	}

	confirmRequest := &trade.ConfirmRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   c.Request.Context(),
			From:  *req.From,
			To:    *req.To,
			Nonce: *req.Nonce,
		},
	}

	if err := h.Service.ManualConfirm(confirmRequest); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errcode.ErrOK.String(),
		"message": "OK",
	})
}
