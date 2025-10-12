// internal/middleware/hmac_simple.go
package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HMACConfig struct {
	Secret    string
	SkipPaths []string
}

const hdrSignature = "X-Signature"

func RequireHMAC(cfg HMACConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		for _, s := range cfg.SkipPaths {
			if strings.HasPrefix(p, s) {
				c.Next()
				return
			}
		}
		if strings.TrimSpace(cfg.Secret) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server hmac not configured"})
			return
		}

		sigB64 := strings.TrimSpace(c.GetHeader(hdrSignature))
		if sigB64 == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing HMAC header"})
			return
		}

		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		if !verifyHMAC([]byte(cfg.Secret), rawBody, sigB64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		c.Next()
	}
}

func verifyHMAC(secret, payload []byte, sigB64 string) bool {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	sum := mac.Sum(nil)

	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return false
	}
	return hmac.Equal(sum, sig)
}
