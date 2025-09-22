package handler

import (
	"basic-gin/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/:id", h.GetByID)            // GET    /accounts/:id
	rg.GET("", h.ListByClient)           // GET    /accounts?client_id=123
	rg.POST("", h.Create)                // POST   /accounts
	rg.POST("/:id/deposit", h.Deposit)   // POST   /accounts/:id/deposit
	rg.POST("/:id/withdraw", h.Withdraw) // POST  /accounts/:id/withdraw
}

type accountCreateReq struct {
	ClientID int `json:"client_id" binding:"required,min=1"`
}

type amountReq struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *AccountHandler) GetByID(c *gin.Context) {
	id, err := parseInt(c.Param("id"))
	if err != nil || id <= 0 {
		h.respondError(c, http.StatusBadRequest, errOr("invalid id", err))
		return
	}
	ctx := c.Request.Context()
	res, err := h.svc.GetById(ctx, id)
	if err != nil {
		code := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			code = http.StatusNotFound
		}
		h.respondError(c, code, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) ListByClient(c *gin.Context) {
	q := c.Query("client_id")
	if q == "" {
		h.respondError(c, http.StatusBadRequest, errStr("missing query param: client_id"))
		return
	}
	clientID, err := parseInt(q)
	if err != nil || clientID <= 0 {
		h.respondError(c, http.StatusBadRequest, errOr("invalid client_id", err))
		return
	}
	ctx := c.Request.Context()
	res, err := h.svc.GetByClientId(ctx, clientID)
	if err != nil {
		code := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			code = http.StatusNotFound
		}
		h.respondError(c, code, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) Create(c *gin.Context) {
	var in accountCreateReq
	if err := c.ShouldBindJSON(&in); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	out, err := h.svc.Save(ctx, in.ClientID)
	if err != nil {
		code := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			code = http.StatusNotFound
		}
		h.respondError(c, code, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *AccountHandler) Deposit(c *gin.Context) {
	id, err := parseInt(c.Param("id"))
	if err != nil || id <= 0 {
		h.respondError(c, http.StatusBadRequest, errOr("invalid id", err))
		return
	}
	var in amountReq
	if err := c.ShouldBindJSON(&in); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	out, err := h.svc.Deposit(ctx, id, in.Amount)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *AccountHandler) Withdraw(c *gin.Context) {
	id, err := parseInt(c.Param("id"))
	if err != nil || id <= 0 {
		h.respondError(c, http.StatusBadRequest, errOr("invalid id", err))
		return
	}
	var in amountReq
	if err := c.ShouldBindJSON(&in); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
	out, err := h.svc.Withdraw(ctx, id, in.Amount)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *AccountHandler) respondError(c *gin.Context, code int, err error) {
	rid := c.Writer.Header().Get("X-Request-ID")
	c.JSON(code, gin.H{
		"error":       err.Error(),
		"status_code": code,
		"request_id":  rid,
	})
}

func parseInt(s string) (int, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	return int(i64), err
}

func errStr(msg string) error { return &simpleError{msg} }
func errOr(msg string, base error) error {
	if base == nil {
		return errStr(msg)
	}
	return base
}

type simpleError struct{ s string }

func (e *simpleError) Error() string { return e.s }
