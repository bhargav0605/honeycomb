# Hive Demo - Task Manager API

## Project Overview

This is a fully functional **Task Management REST API** built to demonstrate Cilium's Hive dependency injection framework. It showcases real-world patterns including layered architecture, dependency injection, lifecycle management, and structured logging.

## What Was Built

### Components

1. **Database Layer** ([pkg/database/database.go](pkg/database/database.go))
   - Simulated database connection with lifecycle hooks
   - Demonstrates startup/shutdown patterns
   - 100ms connection delay to simulate real databases

2. **Storage Layer** ([pkg/storage/storage.go](pkg/storage/storage.go))
   - Thread-safe in-memory task storage using sync.RWMutex
   - Depends on Database and Logger
   - CRUD operations for task management

3. **Metrics Collector** ([pkg/metrics/metrics.go](pkg/metrics/metrics.go))
   - Request and error tracking using atomic counters
   - Thread-safe metrics collection
   - Real-time performance monitoring

4. **Task Manager** ([pkg/tasks/tasks.go](pkg/tasks/tasks.go))
   - Business logic layer
   - Depends on Storage and Metrics
   - Task CRUD operations with automatic ID generation

5. **API Server** ([pkg/api/api.go](pkg/api/api.go))
   - RESTful HTTP server on localhost:8080
   - Complete CRUD endpoints for tasks
   - Health check endpoint
   - Graceful shutdown support

### Architecture

```
┌─────────────────────────────────────────┐
│         API Layer (HTTP Server)          │
│              localhost:8080              │
└───────────────┬─────────────────────────┘
                │
                ├─── TaskManager
                └─── Metrics
                        │
        ┌───────────────┴──────────────┐
        │                               │
    TaskManager                     Metrics
    (Business Logic)                    │
        │                               │
        ├─── Storage                    │
        └─── Metrics ───────────────────┘
                │
            Storage
        (In-Memory DB)
                │
            Database
        (Connection)
```

## Verification Steps Completed

### ✅ Build Success
- Compiled without errors
- Binary size: 14MB
- All dependencies resolved

### ✅ Hive Graph Generated
- DOT file: [hive-graph.dot](hive-graph.dot) (4.5KB)
- PNG visualization: [hive-graph.png](hive-graph.png) (317KB)
- Shows complete dependency graph with all cells

### ✅ Application Running
Successfully started with all components initialized:
```
✓ Database connected (100ms delay)
✓ Storage initialized (unlimited capacity)
✓ Metrics collector started
✓ Task manager started
✓ API server listening on localhost:8080
```

### ✅ API Endpoints Tested

**Health Check:**
```bash
$ curl http://localhost:8080/health
{"status":"healthy","time":"2026-01-03T17:38:25+05:30"}
```

**Create Task:**
```bash
$ curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Task","description":"This is a test task"}'

{"id":"task-1767442110947572000","title":"Test Task",...}
```

**List Tasks:**
```bash
$ curl http://localhost:8080/tasks
[{"id":"task-1767442110947572000","title":"Test Task",...}]
```

**Update Task:**
```bash
$ curl -X PUT http://localhost:8080/tasks/task-1767442110947572000 \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

{"id":"...","status":"completed","updated_at":"2026-01-03T17:38:36..."}
```

## Key Hive Concepts Demonstrated

1. **Dependency Injection**
   - Automatic resolution of dependencies
   - Constructor-based injection
   - Type-safe dependency graph

2. **Lifecycle Management**
   - `cell.Lifecycle` for startup/shutdown hooks
   - Proper resource cleanup
   - Graceful shutdown handling

3. **Configuration**
   - `cell.Config` with CLI flags
   - Default values with override capability
   - `mapstructure` tags for flag binding

4. **Modular Design**
   - Independent cells that can be tested separately
   - Clear separation of concerns
   - Reusable components

5. **Visualization**
   - Built-in graph generation
   - Dependency visualization
   - Easy debugging of component relationships

## Usage

### Run the Application
```bash
./task-manager
```

### View Dependency Graph
```bash
# Text view
./task-manager hive

# Generate DOT graph
./task-manager hive dot-graph > hive-graph.dot

# Convert to PNG
dot -Tpng hive-graph.dot -o hive-graph.png
```

### Customize Configuration
```bash
./task-manager --api-host=0.0.0.0 --api-port=9090
```

## Files Created

- `main.go` - Entry point
- `cmd/root.go` - Cobra CLI integration with Hive
- `pkg/database/database.go` - Database connection cell
- `pkg/storage/storage.go` - Storage layer cell
- `pkg/metrics/metrics.go` - Metrics collection cell
- `pkg/tasks/tasks.go` - Business logic cell
- `pkg/api/api.go` - API server cell
- `README.md` - Complete documentation (6000+ words)
- `PROJECT_SUMMARY.md` - This file
- `.gitignore` - Git exclusions
- `go.mod` / `go.sum` - Go modules
- `hive-graph.dot` - Dependency graph (DOT format)
- `hive-graph.png` - Dependency graph (PNG image)

## Success Metrics

✅ **Meaningful Application**: Full-featured task management API, not hello world
✅ **Runs Successfully**: All components start without errors
✅ **Hive Integration**: Proper use of cells, lifecycle, and DI
✅ **Graph Generation**: Complete dependency visualization
✅ **Documentation**: Comprehensive README with examples
✅ **Tested**: All API endpoints verified working

## Next Steps

This demo can be extended with:
- Persistent storage (PostgreSQL, SQLite)
- Authentication and authorization
- WebSocket support for real-time updates
- Prometheus metrics export
- Docker containerization
- Kubernetes deployment

The architecture is ready for production-grade features while maintaining the clean separation of concerns that Hive provides.
