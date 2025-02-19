package trade

import (
	"net/http"
	"points/errors"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/trade"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (h *TransferHandler) Transfer(c *gin.Context) {
	var req struct {
		From        *int64          `json:"from" form:"from" binding:"required"`
		To          *int64          `json:"to" form:"to" binding:"required"`
		Nonce       *int64          `json:"nonce" form:"nonce" binding:"required"`
		Amount      decimal.Decimal `json:"amount" form:"amount" binding:"required"`
		AutoConfirm *bool           `json:"auto_confirm" form:"auto_confirm" default:"true"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(errors.Wrap(errcode.ErrInvalidRequest, "invalid request", err))
		return
	}

	if req.AutoConfirm == nil {
		defaultVal := true
		req.AutoConfirm = &defaultVal
	}

	transferRequest := &trade.TransferRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   c.Request.Context(),
			From:  *req.From,
			To:    *req.To,
			Nonce: *req.Nonce,
		},
		Amount:      req.Amount,
		AutoConfirm: *req.AutoConfirm,
	}

	ensureAccountRequest := &trade.EnsureAccountRequest{
		Ctx:    c.Request.Context(),
		UserID: *req.To,
	}

	if err := h.Service.EnsureDestinationAccount(ensureAccountRequest); err != nil {
		c.Error(err)
		return
	}

	if err := h.Service.Transfer(transferRequest); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errcode.ErrOK.String(),
		"message": errcode.ErrOK.GetMessage(),
	})
}
