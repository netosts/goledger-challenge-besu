package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netosts/goledger-challenge-besu/internal/models"
	"github.com/netosts/goledger-challenge-besu/internal/usecases"
)

type Handler struct {
	contractUseCase *usecases.ContractUseCase
}

func NewHandler(contractUseCase *usecases.ContractUseCase) *Handler {
	return &Handler{
		contractUseCase: contractUseCase,
	}
}

func (h *Handler) SetValue(c *gin.Context) {
	var req models.SetValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
		return
	}

	if err := h.contractUseCase.SetValue(req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to set value: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Value set successfully",
	})
}

func (h *Handler) GetValue(c *gin.Context) {
	value, err := h.contractUseCase.GetValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get value: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ValueResponse{
		Value: value,
	})
}

func (h *Handler) SyncValue(c *gin.Context) {
	if err := h.contractUseCase.SyncValue(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to sync value: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Value synchronized successfully",
	})
}

func (h *Handler) CheckValue(c *gin.Context) {
	result, err := h.contractUseCase.CheckValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to check values: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Service is healthy",
	})
}
