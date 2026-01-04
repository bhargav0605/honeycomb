package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/cilium/hive/cell"
)

// Cell provides database connection management
var Cell = cell.Module(
	"database",
	"Database Connection Manager",

	cell.Provide(newDatabase),
)

// Database represents a database connection (simulated)
type Database interface {
	Ping(ctx context.Context) error
	IsConnected() bool
}

type db struct {
	logger    *slog.Logger
	connected bool
}

// newDatabase creates a new database connection with lifecycle hooks
func newDatabase(lc cell.Lifecycle, logger *slog.Logger) Database {
	d := &db{
		logger:    logger.With("component", "database"),
		connected: false,
	}

	lc.Append(cell.Hook{
		OnStart: func(ctx cell.HookContext) error {
			d.logger.Info("Connecting to database...")
			// Simulate connection time
			time.Sleep(100 * time.Millisecond)
			d.connected = true
			d.logger.Info("Database connected successfully")
			return nil
		},
		OnStop: func(ctx cell.HookContext) error {
			d.logger.Info("Closing database connection...")
			d.connected = false
			d.logger.Info("Database connection closed")
			return nil
		},
	})

	return d
}

func (d *db) Ping(ctx context.Context) error {
	if !d.connected {
		return context.DeadlineExceeded
	}
	return nil
}

func (d *db) IsConnected() bool {
	return d.connected
}
