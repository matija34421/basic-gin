package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")

		if rid == "" {
			rid = uuid.NewString()
		}

		c.Writer.Header().Set("X-Request-ID", rid)
		c.Set("RequestIDKey", rid)

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, rid)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
