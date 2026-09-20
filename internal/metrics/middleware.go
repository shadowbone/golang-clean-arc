package metrics

import (
	"golang-rest-api/internal/apperror"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Middleware(m *Metrics) fiber.Handler {
	return func(c fiber.Ctx) error {
		m.RequestInFlight.Inc()
		defer m.RequestInFlight.Dec()

		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			status = apperror.StatusFor(err)
		}

		route := c.Route().Path
		if route == "" {
			route = "unknown"
		}

		lbl := prometheusLabels(c.Method(), route, status)

		m.RequestDuration.With(lbl).Observe(time.Since(start).Seconds())
		m.RequestTotal.With(lbl).Inc()

		return err
	}
}

func prometheusLabels(method, route string, status int) map[string]string {
	return map[string]string{
		"method": method,
		"route":  route,
		"status": strconv.Itoa(status),
	}
}
