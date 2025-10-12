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

// @Summary      Kreiraj transakciju (transfer)
// @Description  Kreira transfer između dva računa i vraća snimljenu transakciju.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body  body   dto.TransactionCreate       true  "Podaci o transferu"
// @Success      200   {object}  dto.TransactionResponse  "Uspešno kreirano"
// @Failure      400   {object}  map[string]string        "Neispravan JSON ulaz; npr: {\"error\":\"invalid character...\"}"
// @Router       /transactions [post]
// @Security HmacSignature
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

// @Summary      Transakcije po računu
// @Description  Lista transakcija za dati račun uz paginaciju (limit/offset).
// @Tags         transactions
// @Produce      json
// @Param        accountID  path   int  true   "ID računa (pozitivni integer)"
// @Param        limit      query  int  false  "Broj po stranici"  default(50)  minimum(1) maximum(200)
// @Param        offset     query  int  false  "Pomak (offset)"    default(0)   minimum(0)
// @Success      200        {array}  dto.TransactionResponse
// @Failure      400        {object} map[string]string  "Neispravan accountID; npr: {\"error\":\"invalid account id\"}"
// @Router       /transactions/by-account/{accountID} [get]
// @Security	 HmacSignature
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
