# Go SSE Sample with DDD

A Go-based Server-Sent Events (SSE) sample application demonstrating Domain-Driven Design (DDD) architecture. This project provides a RESTful API for managing metrics and metric readings, with real-time event broadcasting through SSE.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [SSE Events](#sse-events)
- [Architecture Details](#architecture-details)
- [API Examples](#api-examples)
- [Development](#development)

## Overview

This application demonstrates how to build a scalable SSE-based event streaming system using Go and DDD principles. It manages metrics (e.g., temperature, CPU usage) and their readings, broadcasting events in real-time to connected clients when new metrics or readings are created.

### Key Concepts

- **Metrics**: Named entities that represent measurable quantities (e.g., "CPU Temperature", "Memory Usage")
- **Metric Readings**: Time-series data points associated with a specific metric
- **SSE Events**: Real-time notifications sent to connected clients when metrics or readings are created

## Architecture

The project follows Domain-Driven Design (DDD) principles with a clean architecture approach:

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                     │
│  (Controllers, DTOs) - HTTP handlers, request/response   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                          │
│  (Entities, Use Cases, Enums) - Business logic          │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  (In-memory implementations) - Data persistence         │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Package Layer                         │
│  (SSE Hub, Event Store) - Infrastructure                 │
└─────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

- **Presentation Layer**: Handles HTTP requests/responses, validates input, converts between DTOs and domain entities
- **Domain Layer**: Contains business logic, entities, and use cases
- **Repository Layer**: Provides data persistence abstractions (currently in-memory implementations)
- **Package Layer**: Infrastructure components (SSE hub, event store)

## Features

- ✅ RESTful API for metric management
- ✅ RESTful API for metric reading creation
- ✅ Real-time event streaming via Server-Sent Events (SSE)
- ✅ Event replay support using `Last-Event-ID` header
- ✅ Automatic event retention with configurable TTL
- ✅ Client connection management with maximum limit (10,000 clients)
- ✅ Graceful server shutdown
- ✅ UUID v7 for time-ordered identifiers
- ✅ Input validation and error handling
- ✅ Thread-safe in-memory repositories

## Technology Stack

### Server
- **Go**: 1.23.5
- **Gin**: HTTP web framework
- **UUID**: v7 for time-ordered unique identifiers
- **In-memory storage**: For metrics, readings, and events

### Client (Example)
- **Node.js**: v22.2.0
- **EventSource**: Node.js SSE client library

## Project Structure

```
go-sse-sample/
├── cmd/
│   ├── server/
│   │   └── main.go              # Application entry point
│   └── client/
│       ├── client.mjs           # Node.js SSE client example
│       ├── package.json
│       └── package-lock.json
├── internal/
│   ├── domain/
│   │   ├── entity/              # Domain entities
│   │   │   ├── metric_entity.go
│   │   │   └── metric_reading_entity.go
│   │   ├── enum/                # Domain enumerations
│   │   │   └── event_types.go
│   │   └── use_case/            # Business logic
│   │       ├── metric_use_case.go
│   │       └── metric_reading_use_case.go
│   ├── presentation/
│   │   ├── controller/          # HTTP controllers
│   │   │   ├── metric_http_gin_controller.go
│   │   │   ├── metric_reading_http_gin_controller.go
│   │   │   └── events_http_gin_controller.go
│   │   └── dto/                 # Data Transfer Objects
│   │       ├── metric_dto.go
│   │       ├── metric_reading_dto.go
│   │       └── event_dto.go
│   └── repository/              # Data persistence
│       ├── metric_inmemory.go
│       ├── metric_reading_inmemory.go
│       └── event_store_inmemory.go
├── pkg/
│   └── sse/                     # SSE infrastructure
│       ├── sse_hub.go           # SSE hub for client management
│       ├── client.go            # SSE client implementation
│       ├── event.go             # Event structure
│       └── event_store.go      # Event store interface
├── docs/
│   └── api/                     # API documentation
│       ├── metrics_api_docs.http
│       └── metric_readings_api_docs.http
├── .nvmrc                       # Node.js version specification
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23.5 or later
- Node.js v22.2.0 or later (for running the example client)
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-sse-sample
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install Node.js client dependencies (optional, for running the example client):
```bash
cd cmd/client
npm install
cd ../..
```

### Running the Server

Using Make:
```bash
make run-server
```

Or directly with Go:
```bash
go run cmd/server/main.go
```

The server will start on port `8089`. You should see:
```
server started on port 8089
```

### Running the Example Client

A Node.js example client is provided to demonstrate SSE connectivity. Make sure the server is running first:

```bash
make run-client
# or
node cmd/client/client.mjs
```

### Stopping the Server

Press `Ctrl+C` to initiate graceful shutdown. The server will:
1. Stop accepting new connections
2. Stop event retention
3. Wait up to 1 minute for existing connections to close
4. Shut down gracefully

## API Examples

For detailed API documentation and examples, see the `docs/api/` directory:

- `docs/api/metrics_api_docs.http` - Metrics API examples
- `docs/api/metric_readings_api_docs.http` - Metric readings API examples

These HTTP files can be used with REST Client extensions (VS Code, IntelliJ) or tools like Postman, cURL, or HTTPie.

### Available Endpoints

- `POST /metrics` - Create a new metric
- `GET /metrics/:id` - Get metric by ID
- `POST /metrics/readings` - Create a metric reading
- `GET /events/watch` - SSE endpoint for real-time events

A Node.js client example is available in `cmd/client/client.mjs`.

## SSE Events

The application broadcasts real-time events via Server-Sent Events (SSE) when metrics or readings are created.

### Event Types

- **`metric_created`**: Emitted when a new metric is created
- **`metric_reading_created`**: Emitted when a new metric reading is created

### Key Features

- **Event Replay**: Clients can reconnect using the `Last-Event-ID` header to receive missed events
- **Event Retention**: Events are stored in memory with a configurable TTL (default: 1 minute)
- **Connection Management**: Supports up to 10,000 concurrent SSE clients

See `cmd/client/client.mjs` for a complete Node.js client example.

## Architecture Details

### Domain Entities

#### Metric
- **ID**: UUID v7 (time-ordered)
- **Name**: String (required, non-empty)
- **Validation**: ID must be UUID v7, name cannot be empty

#### MetricReading
- **ID**: UUID v7 (time-ordered)
- **MetricID**: UUID v7 reference to Metric
- **Value**: Float64 (must be > 0)
- **Timestamp**: Time (UTC, defaults to current time if not provided)
- **Validation**: All fields validated, UUIDs must be v7

### Use Cases

#### MetricUseCase
- `CreateMetric`: Creates a new metric and broadcasts `metric_created` event
- `GetMetricByID`: Retrieves a metric by ID

#### MetricReadingUseCase
- `CreateMetricReading`: Creates a new reading, validates metric exists, and broadcasts `metric_reading_created` event

### SSE Hub

The SSE Hub manages client connections and event broadcasting:

- **Client Registration**: Registers new SSE clients
- **Client Unregistration**: Removes disconnected clients
- **Event Broadcasting**: Broadcasts events to all connected clients
- **Client Limit**: Maximum 10,000 concurrent clients (oldest disconnected if limit reached)
- **Slow Client Handling**: Drops clients that cannot keep up with event stream

### Event Store

The in-memory event store:
- Stores events for replay functionality
- Implements TTL-based retention
- Provides `GetEventsAfterID` for event replay
- Thread-safe operations

## Configuration

Default configuration (in `cmd/server/main.go`):
- **Port**: `8089`
- **Max SSE Clients**: `10,000`
- **Event TTL**: `1 minute`
- **Graceful Shutdown Timeout**: `1 minute`

To modify these values, edit the constants and variables in `main.go`.

## Development

### Building

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

### Testing

The project follows Go best practices with interface-based design for testability. Use the provided HTTP files in `docs/api/` for API testing with REST Client extensions or tools like Postman, cURL, or HTTPie.

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
