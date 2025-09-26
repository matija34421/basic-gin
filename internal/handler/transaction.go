package handler

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: *s}
}

func (h *TransactionHandler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("/by-account/:accountID", h.ListByAccountID)
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var in dto.TransactionCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.transactionService.CreateTransfer(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *TransactionHandler) ListByAccountID(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil || accountID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	out, err := h.transactionService.ListByAccountID(c.Request.Context(), accountID, limit, offset)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
