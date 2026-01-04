# Task Manager API - Hive Demo

A practical demonstration of [Cilium's Hive](https://github.com/cilium/hive) dependency injection framework in Go.

This is a fully functional Task Management REST API that showcases Hive's core concepts:
- Modular cell-based architecture
- Dependency injection
- Lifecycle management
- Configuration management
- Visualization tools

## 🎯 What This Demonstrates

Unlike simple "Hello World" examples, this project shows real-world patterns:

- **Multi-layered Architecture**: Infrastructure → Business Logic → API
- **Dependency Management**: Components explicitly declare dependencies
- **Lifecycle Hooks**: Proper startup/shutdown ordering
- **Configuration**: Flag-based configuration with defaults
- **Metrics Collection**: Request/error tracking
- **Thread Safety**: Concurrent access patterns
- **RESTful API**: Complete CRUD operations

## 🏗️ Architecture

```
Task Manager
├── Infrastructure Layer
│   ├── Logger (structured logging)
│   ├── Database (connection management)
│   ├── Storage (in-memory store)
│   └── Metrics (request tracking)
├── Business Logic Layer
│   └── Tasks (task management logic)
└── API Layer
    └── HTTP Server (REST API)
```

### Dependency Graph

```
API Server
├── TaskManager
│   ├── Storage
│   │   └── Database
│   └── Metrics
└── Logger

All components depend on Logger
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- (Optional) Graphviz for dependency visualization

### Installation

```bash
# Clone the repository
cd /path/to/hive-demo

# Download dependencies
go mod download

# Build the application
go build -o task-manager

# Run the application
./task-manager
```

### Running the API

```bash
# Start with default settings (port 8080)
./task-manager

# Start with custom port
./task-manager --api-port 3000

# Start with debug logging
./task-manager --log-level debug

# Start with custom host
./task-manager --api-host 0.0.0.0 --api-port 8080
```

The API will be available at `http://localhost:8080`

## 📊 Visualizing the Hive Architecture

### Method 1: Text View

See all components and their dependencies:

```bash
./task-manager hive
```

Output shows:
- All registered cells (modules)
- Their configurations
- Constructor dependencies (⇨ inputs, ⇦ outputs)
- Start/Stop hooks in order

### Method 2: Dependency Graph (Visual)

Generate a visual dependency graph:

```bash
# Generate DOT file
./task-manager hive dot-graph > hive-graph.dot

# Convert to PNG (requires graphviz)
dot -Tpng hive-graph.dot -o hive-graph.png

# Convert to SVG (better for web)
dot -Tsvg hive-graph.dot -o hive-graph.svg

# Open the graph
open hive-graph.png  # macOS
xdg-open hive-graph.png  # Linux
```

#### Install Graphviz

```bash
# macOS
brew install graphviz

# Ubuntu/Debian
sudo apt-get install graphviz

# Windows (via Chocolatey)
choco install graphviz
```

### Method 3: Interactive Inspection

```bash
# Show help for hive commands
./task-manager hive --help

# Options:
#   hive              Show all cells
#   hive dot-graph    Generate dependency graph
```

## 🔌 API Endpoints

### Root
```bash
GET http://localhost:8080/
```
Returns API information and available endpoints.

### Health Check
```bash
GET http://localhost:8080/health
```
Returns service health status.

### Statistics
```bash
GET http://localhost:8080/stats
```
Returns metrics (total tasks, requests, errors, status breakdown).

### List Tasks
```bash
GET http://localhost:8080/tasks
```

### Create Task
```bash
POST http://localhost:8080/tasks
Content-Type: application/json

{
  "title": "Learn Hive",
  "description": "Study Cilium's dependency injection framework"
}
```

### Get Task
```bash
GET http://localhost:8080/tasks/{task-id}
```

### Update Task
```bash
PUT http://localhost:8080/tasks/{task-id}
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed"
}
```

### Delete Task
```bash
DELETE http://localhost:8080/tasks/{task-id}
```

## 🧪 Testing the API

### Using curl

```bash
# Get API info
curl http://localhost:8080/

# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Task","description":"Test the API"}'

# List all tasks
curl http://localhost:8080/tasks

# Get stats
curl http://localhost:8080/stats

# Update a task (replace task-id with actual ID)
curl -X PUT http://localhost:8080/tasks/task-1234567890 \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

# Delete a task
curl -X DELETE http://localhost:8080/tasks/task-1234567890
```

### Using httpie

```bash
# Install httpie
brew install httpie  # macOS
sudo apt-get install httpie  # Ubuntu

# Create a task
http POST :8080/tasks title="Learn Hive" description="Complete tutorial"

# List tasks
http :8080/tasks

# Update task
http PUT :8080/tasks/task-123 status=completed

# Delete task
http DELETE :8080/tasks/task-123
```

## 📁 Project Structure

```
hive-demo/
├── main.go                 # Application entry point
├── cmd/
│   └── root.go            # CLI command setup & Hive initialization
├── pkg/
│   ├── api/
│   │   └── api.go         # HTTP API server (depends on tasks, metrics)
│   ├── database/
│   │   └── database.go    # Database connection (simulated)
│   ├── logger/
│   │   └── logger.go      # Structured logging
│   ├── metrics/
│   │   └── metrics.go     # Metrics collection
│   ├── storage/
│   │   └── storage.go     # In-memory storage (depends on database)
│   └── tasks/
│       └── tasks.go       # Task business logic (depends on storage, metrics)
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── README.md              # This file
```

## 🧩 Key Hive Concepts Demonstrated

### 1. Cell Definition

Each component is a cell:

```go
var Cell = cell.Module(
    "component-id",
    "Human-readable description",

    cell.Config(defaultConfig),  // Optional configuration
    cell.Provide(newComponent),  // Constructor
)
```

### 2. Dependency Injection

Dependencies declared in constructor parameters:

```go
func newComponent(
    logger *slog.Logger,      // Injected automatically
    storage Storage,          // Injected automatically
    metrics Metrics,          // Injected automatically
) Component {
    // Use dependencies
}
```

### 3. Lifecycle Management

Components register start/stop hooks:

```go
lc.Append(cell.Hook{
    OnStart: func(ctx context.Context) error {
        // Initialize component
        return nil
    },
    OnStop: func(ctx context.Context) error {
        // Cleanup
        return nil
    },
})
```

### 4. Configuration

Configuration with CLI flags:

```go
type Config struct {
    Port int
}

func (c Config) Flags(flags *pflag.FlagSet) {
    flags.Int("api-port", c.Port, "API server port")
}
```

### 5. Module Composition

Building the application from cells:

```go
var App = cell.Module(
    "app",
    "Application",

    logger.Cell,     // Infrastructure
    database.Cell,   // Infrastructure
    storage.Cell,    // Infrastructure
    metrics.Cell,    // Infrastructure
    tasks.Cell,      // Business logic
    api.Cell,        // API layer
)
```

## 🎓 Learning Path

1. **Start Simple**: Look at `pkg/logger/logger.go` - simplest cell
2. **Add Dependencies**: Check `pkg/storage/storage.go` - depends on database & logger
3. **Business Logic**: Examine `pkg/tasks/tasks.go` - depends on storage & metrics
4. **API Layer**: Study `pkg/api/api.go` - depends on tasks & metrics
5. **Composition**: See `cmd/root.go` - how everything wires together
6. **Visualization**: Run `hive` command to see the result

## 🔍 Understanding the Flow

### Startup Sequence

1. **Hive Construction**: `hive.New(App)` analyzes all cells
2. **Dependency Resolution**: Builds dependency graph
3. **Topological Sort**: Determines initialization order
4. **Object Creation**: Instantiates components in order
5. **Hook Execution**: Runs OnStart hooks sequentially
6. **Ready**: Application running

### Request Flow

```
HTTP Request
    ↓
API Server (logs request, increments metrics)
    ↓
Task Manager (business logic)
    ↓
Storage (thread-safe data access)
    ↓
HTTP Response
```

### Shutdown Sequence

1. **Signal Received**: SIGINT or SIGTERM
2. **Hook Execution**: OnStop hooks in reverse order
3. **API Server**: Stops accepting requests, drains connections
4. **Task Manager**: Reports statistics
5. **Storage**: Clears data
6. **Database**: Closes connections
7. **Clean Exit**: Process terminates

## 🤔 Why Hive?

### Problems Solved

✅ **No manual dependency wiring** - Framework handles it
✅ **No initialization order bugs** - Automatic ordering
✅ **No global variables** - Everything injected
✅ **Easy testing** - Inject mocks trivially
✅ **Self-documenting** - `hive` command shows architecture
✅ **Parallel development** - Teams work on independent cells

### When to Use

- Applications with 10+ components
- Complex dependency graphs
- Multiple developers
- Long-running services
- High testability requirements

### When NOT to Use

- Simple CLI tools (startup time matters)
- < 5 components (overhead not worth it)
- Performance-critical initialization

## 📚 Further Reading

- [Hive Documentation](https://pkg.go.dev/github.com/cilium/hive)
- [Cilium's Guide to the Hive](https://docs.cilium.io/en/stable/contributing/development/hive/)
- [uber-go/dig](https://pkg.go.dev/go.uber.org/dig) - Underlying DI library
- [Cilium GitHub](https://github.com/cilium/cilium) - Production example

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Check what's using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Run on different port
./task-manager --api-port 3000
```

### Dependencies Not Found

```bash
# Clean and rebuild
go clean -modcache
go mod download
go build
```

### Graph Generation Fails

```bash
# Make sure graphviz is installed
dot -V

# If not installed, see "Install Graphviz" section above
```

## 📝 License

MIT License - feel free to use this as a learning resource or starting template.

## 🙏 Credits

- [Cilium Project](https://github.com/cilium/cilium) - For the Hive framework
- [uber-go/dig](https://github.com/uber-go/dig) - Underlying DI library

---

**Happy Learning! 🚀**

If you found this helpful, give it a ⭐ on GitHub!
