package trade

import (
	"net/http"
	"points/pkg/models/enum/errcode"
	"points/pkg/models/trade"
	"points/pkg/module/httputil"

	"github.com/gin-gonic/gin"
)

func (h *TransferHandler) Cancel(c *gin.Context) {
	var req struct {
		From  *int32 `json:"from" form:"from" binding:"required"`
		To    *int32 `json:"to" form:"to" binding:"required"`
		Nonce *int64 `json:"nonce" form:"nonce" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transferRequest := &trade.CancelRequest{
		BaseRequest: trade.BaseRequest{
			Ctx:   c.Request.Context(),
			From:  *req.From,
			To:    *req.To,
			Nonce: *req.Nonce,
		},
	}

	if err := h.Validator.ValidateCancelRequest(transferRequest); err != nil {
		status, resp := httputil.FormatError(err, "validate failed", http.StatusBadRequest)
		c.JSON(status, resp)
		return
	}

	if err := h.Service.Cancel(transferRequest); err != nil {
		status, resp := httputil.FormatError(err, "cancel failed", http.StatusInternalServerError)
		c.JSON(status, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  errcode.ErrOK.String(),
		"message": "OK",
	})
}
