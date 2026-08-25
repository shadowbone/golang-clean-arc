package logger

import (
	"encoding/json"
	"fmt"
	"golang-rest-api/internal/apperror"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

const maxBodyLog = 4096

var sensitiveFields = map[string]bool{
	"password":         true,
	"password_confirm": true,
	"token":            true,
	"refresh_token":    true,
	"email":            true,
}

func Middleware(base *slog.Logger, logBody bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		reqId := requestid.FromContext(c)

		log := base.With(
			slog.String("request_id", reqId),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
		)
		if logBody && shouldLogBody(c) {
			log = log.With(
				slog.Any("body", sanitize(c.Body())),
			)
		}

		c.SetContext(WithContext(c.Context(), log))

		err := c.Next()
		status := c.Response().StatusCode()
		if err != nil {
			status = apperror.StatusFor(err)
		}
		attrs := []any{
			slog.Int("status", status),
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", c.IP()),
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
		}

		switch {
		case status >= 500:
			log.Error("request selesai", attrs...)
		case status >= 400:
			log.Warn("request selesai", attrs...)
		default:
			log.Info("request selesai", attrs...)
		}

		return err
	}
}

func shouldLogBody(c fiber.Ctx) bool {
	switch c.Method() {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch:
	default:
		return false
	}

	ct := string(c.Request().Header.ContentType())
	return strings.Contains(ct, "application/json")
}

func sanitize(body []byte) any {
	if len(body) == 0 {
		return nil
	}

	if len(body) > maxBodyLog {
		return fmt.Sprintf("[body %d bytes, terlalu besar]", len(body))
	}

	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return "[body tidak dapat diparsing]"
	}

	redactMap(m)
	return m
}

func redactMap(m map[string]any) {
	for key, val := range m {
		if sensitiveFields[strings.ToLower(key)] {
			m[key] = "[RAHASIA]"
		}

		switch child := val.(type) {
		case map[string]any:
			redactMap(child)
		case []any:
			for _, item := range child {
				if obj, ok := item.(map[string]any); ok {
					redactMap(obj)
				}
			}
		}
	}
}
