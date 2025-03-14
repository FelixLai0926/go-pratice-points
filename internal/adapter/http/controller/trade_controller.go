package controller

import (
	"net/http"
	"points/internal/adapter/http/dto"
	"points/internal/domain"
	"points/internal/domain/command"
	"points/internal/domain/port"
	"points/internal/shared/apperror"
	"points/internal/shared/errcode"
	"points/internal/shared/mapper"

	"github.com/gin-gonic/gin"
)

type TradeController struct {
	TradeUsecase domain.TradeUsecase
	config       port.Config
}

func NewTradeController(usecase domain.TradeUsecase, config port.Config) *TradeController {
	return &TradeController{
		TradeUsecase: usecase,
		config:       config,
	}
}

func (h *TradeController) Transfer(c *gin.Context) {
	var request dto.TransferRequest

	if err := c.ShouldBind(&request); err != nil {
		c.Error(apperror.Wrap(errcode.ErrInvalidRequest, "invalid request", err))
		return
	}
	h.config.SetDefault(&request)

	cmd, err := mapper.MapStruct[command.TransferCommand](h.config, &request)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.TradeUsecase.Transfer(c, cmd); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse())
}

func (h *TradeController) Confirm(c *gin.Context) {
	var request dto.ConfirmRequest

	if err := c.ShouldBind(&request); err != nil {
		c.Error(apperror.Wrap(errcode.ErrInvalidRequest, "invalid request", err))
		return
	}

	cmd, err := mapper.MapStruct[command.ConfirmCommand](h.config, &request)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.TradeUsecase.ManualConfirm(c, cmd); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse())
}

func (h *TradeController) Cancel(c *gin.Context) {
	var request dto.CancelRequest

	if err := c.ShouldBind(&request); err != nil {
		c.Error(apperror.Wrap(errcode.ErrInvalidRequest, "invalid request", err))
		return
	}

	cmd, err := mapper.MapStruct[command.CancelCommand](h.config, &request)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.TradeUsecase.Cancel(c, cmd); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse())
}
