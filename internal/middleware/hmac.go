package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type HMACConfig struct {
	Secret      string
	AllowedSkew time.Duration
	SkipPaths   []string
}

const (
	hdrSignature = "X-Signature"
	hdrTimestamp = "X-Timestamp"
	sigPrefix    = "v1"
)

func RequireHMAC(cfg HMACConfig) gin.HandlerFunc {
	if cfg.AllowedSkew == 0 {
		cfg.AllowedSkew = 5 * time.Minute
	}
	return func(c *gin.Context) {
		//provera za skipovanje ruta u skip path varijabli
		for _, p := range cfg.SkipPaths {
			if c.FullPath() == p || c.Request.URL.Path == p {
				c.Next()
				return
			}
		}

		//provera HMAC headera iz poziva
		sigB64 := strings.TrimSpace(c.GetHeader(hdrSignature))
		tsStr := strings.TrimSpace(c.GetHeader(hdrTimestamp))
		if sigB64 == "" || tsStr == "" {
			abortJSON(c, http.StatusUnauthorized, "missing HMAC headers")
			return
		}

		tsInt, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			abortJSON(c, http.StatusUnauthorized, "invalid timestamp")
			return
		}
		now := time.Now()
		ts := time.Unix(tsInt, 0)
		if ts.Before(now.Add(-cfg.AllowedSkew)) || ts.After(now.Add(cfg.AllowedSkew)) {
			abortJSON(c, http.StatusUnauthorized, "timestamp outside allowed window")
			return
		}

		//citanje body-a
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			abortJSON(c, http.StatusBadRequest, "failed to read body")
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody)) // restore

		//pravljenje payloada
		payload := buildPayload(tsStr, c.Request.Method, c.Request.URL.Path, rawBody)

		ok := verifyHMAC([]byte(cfg.Secret), payload, sigB64)
		if !ok {
			abortJSON(c, http.StatusUnauthorized, "invalid signature")
			return
		}

		c.Next()
	}
}

func buildPayload(ts, method, path string, body []byte) []byte {
	var b bytes.Buffer
	b.Grow(len(ts) + len(method) + len(path) + len(body) + 10)
	b.WriteString(sigPrefix)
	b.WriteByte(':')
	b.WriteString(ts)
	b.WriteByte(':')
	b.WriteString(strings.ToUpper(method))
	b.WriteByte(':')
	b.WriteString(path)
	b.WriteByte(':')
	b.Write(body)
	return b.Bytes()
}

func verifyHMAC(secret, payload []byte, sigB64 string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	sum := mac.Sum(nil)

	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return false
	}
	return hmac.Equal(sum, sig)
}

func abortJSON(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": msg,
	})
}
