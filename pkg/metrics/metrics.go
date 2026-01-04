package metrics

import (
	"log/slog"
	"sync/atomic"

	"github.com/cilium/hive/cell"
)

// Cell provides metrics collection
var Cell = cell.Module(
	"metrics",
	"Metrics Collector",

	cell.Provide(newMetrics),
)

// Metrics provides basic metrics collection
type Metrics interface {
	IncrementRequests()
	IncrementErrors()
	GetRequests() int64
	GetErrors() int64
}

type metrics struct {
	logger   *slog.Logger
	requests atomic.Int64
	errors   atomic.Int64
}

// newMetrics creates a new metrics collector
func newMetrics(lc cell.Lifecycle, logger *slog.Logger) Metrics {
	m := &metrics{
		logger: logger.With("component", "metrics"),
	}

	lc.Append(cell.Hook{
		OnStart: func(ctx cell.HookContext) error {
			m.logger.Info("Metrics collector started")
			return nil
		},
		OnStop: func(ctx cell.HookContext) error {
			m.logger.Info("Metrics summary",
				"total_requests", m.requests.Load(),
				"total_errors", m.errors.Load(),
			)
			return nil
		},
	})

	return m
}

func (m *metrics) IncrementRequests() {
	m.requests.Add(1)
}

func (m *metrics) IncrementErrors() {
	m.errors.Add(1)
}

func (m *metrics) GetRequests() int64 {
	return m.requests.Load()
}

func (m *metrics) GetErrors() int64 {
	return m.errors.Load()
}
