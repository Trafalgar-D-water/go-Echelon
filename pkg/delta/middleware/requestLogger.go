package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const loggerContextKey = "logger"

func RequestLogger(baseLogger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		requestID := newRequestID()

		c.Header("X-Request-ID", requestID)

		requestLogger := baseLogger.With(
			"component", "delta.http",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		c.Set(loggerContextKey, requestLogger)

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(startedAt)

		fields := []any{
			"status", status,
			"duration_ms", duration.Milliseconds(),
		}

		if userID, exists := c.Get("userId"); exists {
			fields = append(fields, "user_id", userID)
		}

		switch {
		case status >= http.StatusInternalServerError:
			requestLogger.Error("http request completed", fields...)
		case status >= http.StatusBadRequest:
			requestLogger.Warn("http request completed", fields...)
		default:
			requestLogger.Info("http request completed", fields...)
		}
	}
}

func Logger(c *gin.Context) *slog.Logger {
	if value, exists := c.Get(loggerContextKey); exists {
		if logger, ok := value.(*slog.Logger); ok {
			return logger
		}
	}

	return slog.Default()
}

func newRequestID() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}
