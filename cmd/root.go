package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bhargavparmar/hive-demo/pkg/api"
	"github.com/bhargavparmar/hive-demo/pkg/database"
	"github.com/bhargavparmar/hive-demo/pkg/metrics"
	"github.com/bhargavparmar/hive-demo/pkg/storage"
	"github.com/bhargavparmar/hive-demo/pkg/tasks"
	"github.com/cilium/hive"
	"github.com/cilium/hive/cell"
	"github.com/spf13/cobra"
)

var (
	// App is the main Hive application containing all components
	App = cell.Module(
		"task-manager",
		"Task Management API",

		// Infrastructure layer - external dependencies
		// Note: Logger is provided automatically by Hive
		database.Cell,
		storage.Cell,
		metrics.Cell,

		// Business logic layer
		tasks.Cell,

		// API layer
		api.Cell,

		// Invoke ensures the API server is constructed and started
		cell.Invoke(func(api.Server) {}),
	)

	// h is the Hive instance shared between commands
	h = hive.New(App)

	// rootCmd is the main command for the application
	rootCmd = &cobra.Command{
		Use:   "task-manager",
		Short: "A simple task management API built with Hive",
		Long: `Task Manager demonstrates Cilium's Hive dependency injection framework.

It includes:
- RESTful API server
- In-memory task storage
- Metrics collection
- Structured logging
- Proper lifecycle management

All components are wired together using Hive's dependency injection.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create a basic logger for Hive
			log := slog.Default()

			// Run the hive (start all components, wait for interrupt)
			if err := h.Run(log); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

// Execute runs the root command
func Execute() {
	// Register all flags from cells
	h.RegisterFlags(rootCmd.Flags())

	// Add hive inspection command
	rootCmd.AddCommand(h.Command())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
