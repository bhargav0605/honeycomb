package logger

import (
	"log/slog"
	"os"

	"github.com/cilium/hive/cell"
	"github.com/spf13/pflag"
)

// Cell provides the logger to the application
var Cell = cell.Module(
	"logger",
	"Structured Logger",

	cell.Config(defaultConfig),
	cell.Provide(newLogger),
)

// Config holds logger configuration
type Config struct {
	Level string
}

var defaultConfig = Config{
	Level: "info",
}

// Flags implements cell.Flagger
func (c Config) Flags(flags *pflag.FlagSet) {
	flags.String("log-level", c.Level, "Log level (debug, info, warn, error)")
}

// newLogger creates a new structured logger
func newLogger(cfg Config) *slog.Logger {
	var level slog.Level

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	logger.Info("Logger initialized", "level", cfg.Level)

	return logger
}
