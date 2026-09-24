package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 记录请求日志，包含 request_id、method、path、status、latency_ms。
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		logger.Info("http request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency,
		)
	}
}

func newRequestID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
