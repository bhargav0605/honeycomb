package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/bhargavparmar/hive-demo/pkg/metrics"
	"github.com/bhargavparmar/hive-demo/pkg/tasks"
	"github.com/cilium/hive/cell"
	"github.com/spf13/pflag"
)

// Cell provides HTTP API server
var Cell = cell.Module(
	"api",
	"HTTP API Server",

	cell.Config(defaultConfig),
	cell.Provide(newServer),
)

// Config holds API server configuration
type Config struct {
	Port int    `mapstructure:"api-port"`
	Host string `mapstructure:"api-host"`
}

var defaultConfig = Config{
	Port: 8080,
	Host: "localhost",
}

// Flags implements cell.Flagger
func (c Config) Flags(flags *pflag.FlagSet) {
	flags.Int("api-port", c.Port, "API server port")
	flags.String("api-host", c.Host, "API server host")
}

// Server represents the HTTP API server
type Server interface {
	Address() string
}

type server struct {
	cfg         Config
	logger      *slog.Logger
	taskManager tasks.TaskManager
	metrics     metrics.Metrics
	httpServer  *http.Server
}

// newServer creates a new HTTP API server with all dependencies
func newServer(lc cell.Lifecycle, cfg Config, logger *slog.Logger, tm tasks.TaskManager, m metrics.Metrics) Server {
	s := &server{
		cfg:         cfg,
		logger:      logger.With("component", "api-server"),
		taskManager: tm,
		metrics:     m,
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/", s.handleTaskByID)
	mux.HandleFunc("/stats", s.handleStats)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lc.Append(cell.Hook{
		OnStart: func(ctx cell.HookContext) error {
			s.logger.Info("Starting API server", "address", s.httpServer.Addr)

			go func() {
				if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					s.logger.Error("API server error", "error", err)
				}
			}()

			s.logger.Info("API server started successfully", "url", fmt.Sprintf("http://%s", s.httpServer.Addr))
			return nil
		},
		OnStop: func(ctx cell.HookContext) error {
			s.logger.Info("Stopping API server...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
				s.logger.Error("Error shutting down server", "error", err)
				return err
			}

			s.logger.Info("API server stopped")
			return nil
		},
	})

	return s
}

func (s *server) Address() string {
	return s.httpServer.Addr
}

// Middleware for logging requests
func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		s.metrics.IncrementRequests()

		s.logger.Info("Request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		s.logger.Info("Response",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"service": "Task Manager API",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"GET /health":       "Health check",
			"GET /stats":        "Get statistics",
			"GET /tasks":        "List all tasks",
			"POST /tasks":       "Create a new task",
			"GET /tasks/{id}":   "Get a specific task",
			"PUT /tasks/{id}":   "Update a task",
			"DELETE /tasks/{id}": "Delete a task",
		},
	}

	s.jsonResponse(w, http.StatusOK, response)
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	s.jsonResponse(w, http.StatusOK, response)
}

func (s *server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.taskManager.GetStats()
	s.jsonResponse(w, http.StatusOK, stats)
}

func (s *server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks := s.taskManager.List()
		s.jsonResponse(w, http.StatusOK, tasks)

	case http.MethodPost:
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		task, err := s.taskManager.Create(req.Title, req.Description)
		if err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusBadRequest, err.Error())
			return
		}

		s.jsonResponse(w, http.StatusCreated, task)

	default:
		s.jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *server) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		s.jsonError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := s.taskManager.Get(id)
		if err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusNotFound, err.Error())
			return
		}
		s.jsonResponse(w, http.StatusOK, task)

	case http.MethodPut:
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		task, err := s.taskManager.Update(id, req.Title, req.Description, req.Status)
		if err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusNotFound, err.Error())
			return
		}

		s.jsonResponse(w, http.StatusOK, task)

	case http.MethodDelete:
		if err := s.taskManager.Delete(id); err != nil {
			s.metrics.IncrementErrors()
			s.jsonError(w, http.StatusNotFound, err.Error())
			return
		}

		s.jsonResponse(w, http.StatusOK, map[string]string{"message": "Task deleted"})

	default:
		s.jsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (s *server) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *server) jsonError(w http.ResponseWriter, status int, message string) {
	s.jsonResponse(w, status, map[string]string{"error": message})
}
