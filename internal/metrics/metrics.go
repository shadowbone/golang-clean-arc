package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	registry        *prometheus.Registry
	RequestDuration *prometheus.HistogramVec
	RequestTotal    *prometheus.CounterVec
	RequestInFlight prometheus.Gauge
	DBPoolTotal     prometheus.GaugeFunc
	DBPoolIdle      prometheus.GaugeFunc
	DBPoolAcquired  prometheus.GaugeFunc
}

type PoolStatsFunc func() (total, idle, acquired int32)

func New(serviceName string, poolStats PoolStatsFunc) *Metrics {
	req := prometheus.NewRegistry()

	labels := prometheus.Labels{"service": serviceName}

	m := &Metrics{
		registry: req,
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "http_request_duration_seconds",
				Help:        "Durasi req HTTP dalam detik",
				ConstLabels: labels,
				Buckets: []float64{
					0.001, 0.005, 0.01, 0.025, 0.05,
					0.1, 0.25, 0.5, 1, 2.5, 5,
				},
			},
			[]string{"method", "route", "status"},
		),
		RequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Jumlah request HTTP",
			},
			[]string{"method", "route", "status"},
		),
		RequestInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:        "http_request_in_flight",
				Help:        "Jumlah request yang sedang diproses",
				ConstLabels: labels,
			},
		),
	}

	req.MustRegister(m.RequestDuration, m.RequestTotal, m.RequestInFlight)

	req.MustRegister(collectors.NewGoCollector())
	req.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	if poolStats != nil {
		m.DBPoolTotal = prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "db_pool_connections_total",
				Help:        "Total koneksi di pool",
				ConstLabels: labels,
			},
			func() float64 { t, _, _ := poolStats(); return float64(t) },
		)

		m.DBPoolIdle = prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "db_pool_connections_idle",
				Help:        "Koneksi idle di pool",
				ConstLabels: labels,
			},
			func() float64 { _, i, _ := poolStats(); return float64(i) },
		)

		m.DBPoolAcquired = prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name:        "db_pool_connections_acquired",
				Help:        "Koneksi yang sedang dipakai",
				ConstLabels: labels,
			},
			func() float64 { _, _, a := poolStats(); return float64(a) },
		)

		req.MustRegister(m.DBPoolTotal, m.DBPoolIdle, m.DBPoolAcquired)
	}

	return m
}

func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}
