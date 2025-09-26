package server

import (
	"log"
	"net/http"
	"time"

	"basic-gin/internal/handler"
	"basic-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newRouter(h *handler.Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "basic-gin-service up",
			"ts":      time.Now().UTC(),
		})
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// clients
	if h == nil || h.ClientHandler == nil {
		log.Println("WARN: client handler is nil - routes will be missing")
	} else {
		clients := v1.Group("/clients")
		h.ClientHandler.Register(clients)
	}

	// accounts
	if h == nil || h.AccountHandler == nil {
		log.Println("WARN: account handler is nil - routes will be missing")
	} else {
		accounts := v1.Group("/accounts")
		h.AccountHandler.Register(accounts)
	}

	//transactions
	if h == nil || h.TransactionHandler == nil {
		log.Println("WARN: client handler is nil - routes will be missing")
	} else {
		transactions := v1.Group("/transactions")
		h.TransactionHandler.Register(transactions)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "route not found",
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"request_id": c.Writer.Header().Get("X-Request-ID"),
		})
	})

	for _, rt := range r.Routes() {
		log.Printf("route: %-6s %-28s -> %s", rt.Method, rt.Path, rt.Handler)
	}

	return r
}
