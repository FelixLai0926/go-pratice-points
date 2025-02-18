package trade

import (
	"net/http"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/trade"
	"points/pkg/module/httputil"

	"github.com/gin-gonic/gin"
)

func (h *TransferHandler) Transfer(c *gin.Context) {
	var req struct {
		From        *int32  `json:"from" form:"from" binding:"required"`
		To          *int32  `json:"to" form:"to" binding:"required"`
		Nonce       *int64  `json:"nonce" form:"nonce" binding:"required"`
		Amount      float64 `json:"amount" form:"amount" binding:"required"`
		AutoConfirm *bool   `json:"auto_confirm" form:"auto_confirm" default:"true"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := h.Validator.ValidateTransferRequest(transferRequest); err != nil {
		status, resp := httputil.FormatError(err, "validate failed", http.StatusBadRequest)
		c.JSON(status, resp)
		return
	}

	EnsureAccountRequest := &trade.EnsureAccountRequest{
		Ctx:    c.Request.Context(),
		UserID: *req.To,
	}

	if err := h.Service.EnsureDestinationAccount(EnsureAccountRequest); err != nil {
		status, resp := httputil.FormatError(err, "ensure destination account failed", http.StatusInternalServerError)
		c.JSON(status, resp)
		return
	}

	if err := h.Service.Transfer(transferRequest); err != nil {
		status, resp := httputil.FormatError(err, "transfer failed", http.StatusInternalServerError)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errcode.ErrOK.String(),
		"message": errcode.ErrOK.GetMessage(),
	})
}
