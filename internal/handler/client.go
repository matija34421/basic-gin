package handler

import (
	"basic-gin/internal/dto"
	"basic-gin/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	svc *service.ClientService
}

func NewClientHandler(svc *service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

func (h *ClientHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.List)        // GET    /clients
	rg.GET("/:id", h.GetByID) // GET    /clients/:id
	rg.POST("", h.Create)     // POST   /clients
	rg.PUT("/:id", h.Update)  // PUT    /clients/:id
}

func (h *ClientHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.svc.GetAll(ctx)
	if err != nil {
		h.respondError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ClientHandler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	res, err := h.svc.GetById(ctx, id)
	if err != nil {
		h.respondError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *ClientHandler) Create(c *gin.Context) {
	var in dto.ClientCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	out, err := h.svc.Save(ctx, in)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *ClientHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	var in dto.ClientUpdate
	if err := c.ShouldBindJSON(&in); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	in.ID = id

	ctx := c.Request.Context()
	out, err := h.svc.Update(ctx, in)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ClientHandler) respondError(c *gin.Context, code int, err error) {
	rid := c.Writer.Header().Get("X-Request-ID")
	c.JSON(code, gin.H{
		"error":       err.Error(),
		"status_code": code,
		"request_id":  rid,
	})
}

func parseID(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }
