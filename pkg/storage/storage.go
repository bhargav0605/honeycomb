package storage

import (
	"log/slog"
	"sync"

	"github.com/bhargavparmar/hive-demo/pkg/database"
	"github.com/cilium/hive/cell"
)

// Cell provides in-memory storage
var Cell = cell.Module(
	"storage",
	"In-Memory Storage",

	cell.Provide(newStorage),
)

// Storage provides thread-safe in-memory storage
type Storage interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	List() map[string]interface{}
	Count() int
}

type memoryStorage struct {
	logger *slog.Logger
	db     database.Database
	mu     sync.RWMutex
	data   map[string]interface{}
}

// newStorage creates a new in-memory storage with database dependency
func newStorage(lc cell.Lifecycle, logger *slog.Logger, db database.Database) Storage {
	s := &memoryStorage{
		logger: logger.With("component", "storage"),
		db:     db,
		data:   make(map[string]interface{}),
	}

	lc.Append(cell.Hook{
		OnStart: func(ctx cell.HookContext) error {
			s.logger.Info("Initializing storage...")
			// Verify database is ready
			if !s.db.IsConnected() {
				s.logger.Warn("Database not connected, storage may have limited functionality")
			}
			s.logger.Info("Storage initialized", "capacity", "unlimited")
			return nil
		},
		OnStop: func(ctx cell.HookContext) error {
			s.logger.Info("Clearing storage...")
			s.mu.Lock()
			defer s.mu.Unlock()
			count := len(s.data)
			s.data = make(map[string]interface{})
			s.logger.Info("Storage cleared", "items_removed", count)
			return nil
		},
	})

	return s
}

func (s *memoryStorage) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	s.logger.Debug("Item stored", "key", key)
}

func (s *memoryStorage) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *memoryStorage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	s.logger.Debug("Item deleted", "key", key)
}

func (s *memoryStorage) List() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]interface{}, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

func (s *memoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
