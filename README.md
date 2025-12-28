# Go SSE Sample with DDD

A Go-based Server-Sent Events (SSE) implementation demonstrating real-time event streaming with Domain-Driven Design (DDD) architecture. Uses a metrics domain as a sample use case to showcase SSE capabilities.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Architecture Details](#architecture-details)
- [API Examples](#api-examples)
- [Demonstration](#demonstration)
- [Development](#development)

## Overview

This application demonstrates how to build a scalable Server-Sent Events (SSE) system in Go using DDD principles. The core focus is on the SSE infrastructure: event broadcasting, client management, event replay, and connection handling.

A metrics domain (metrics and readings) is used as a sample use case to demonstrate SSE functionality. The SSE implementation is domain-agnostic and can be adapted to any event-driven use case.

### Key Concepts

- **SSE Hub**: Manages client connections and broadcasts events to all connected clients
- **Event Store**: Stores events for replay functionality with configurable retention
- **Event Replay**: Clients can reconnect and receive missed events using `Last-Event-ID`
- **Connection Management**: Handles client registration, unregistration, and slow client detection

## Architecture

The project follows Domain-Driven Design (DDD) principles with a clean architecture approach. The SSE infrastructure (`pkg/sse/`) is the core component and is domain-agnostic:

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                     │
│  (Controllers, DTOs) - HTTP handlers, SSE endpoint      │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Domain Layer                          │
│  (Sample domain: Metrics) - Business logic              │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                  Repository Layer                        │
│  (In-memory implementations) - Data persistence         │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    Package Layer (Core)                  │
│  (SSE Hub, Event Store) - SSE Infrastructure             │
│  • Client connection management                          │
│  • Event broadcasting                                    │
│  • Event replay support                                  │
│  • Thread-safe operations                                │
└─────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

- **Presentation Layer**: HTTP handlers, SSE endpoint (`/events/watch`), request/response handling
- **Domain Layer**: Sample domain logic (metrics) that triggers SSE events
- **Repository Layer**: Data persistence abstractions (currently in-memory)
- **Package Layer (Core)**: SSE infrastructure - hub, event store, client management

## Features

### SSE Core Features

- ✅ **Real-time event streaming** via Server-Sent Events (SSE)
- ✅ **Event replay** support using `Last-Event-ID` header
- ✅ **Automatic event retention** with configurable TTL
- ✅ **Client connection management** with maximum limit (10,000 clients)
- ✅ **Slow client detection** - automatically drops clients that can't keep up
- ✅ **Thread-safe operations** for concurrent client handling
- ✅ **Graceful server shutdown** with connection cleanup

### Sample Domain (Metrics)

- RESTful API for creating and retrieving metrics (demonstrates event triggering)
- RESTful API for creating metric readings (demonstrates event broadcasting)

## Technology Stack

### Server
- **Go**: 1.23.5
- **Gin**: HTTP web framework
- **UUID**: v7 for time-ordered unique identifiers
- **In-memory storage**: For metrics, readings, and events

### Client (Demonstration Tool)
- **React**: UI library for building the dashboard
- **Vite**: Fast build tool and dev server
- **Recharts**: Charting library for real-time data visualization
- **EventSource API**: Browser API for SSE connections

## Project Structure

```
go-sse-sample/
├── cmd/
│   ├── server/
│   │   └── main.go              # Application entry point
│   └── client/
│       ├── src/                 # React dashboard (SSE demonstration tool)
│       │   ├── components/      # React components
│       │   ├── hooks/           # Custom hooks (useSSE)
│       │   └── main.jsx         # React entry point
│       ├── index.html
│       ├── package.json
│       └── vite.config.js
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
- Node.js v22.2.0 or later (for running the React dashboard demonstration)
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

3. Install React dashboard dependencies (optional, for running the demonstration dashboard):
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

### Running the React Dashboard (Demonstration Tool)

A React dashboard is provided to visually demonstrate SSE functionality in real-time. This dashboard was created by **Auto (Cursor AI agent)** as a demonstration tool. Make sure the server is running first:

```bash
make run-client
# or
cd cmd/client && npm run dev
```

The dashboard will be available at `http://localhost:5173` (or the port shown in the terminal). It provides:
- Real-time metrics visualization with live charts
- Connection status monitoring
- Debug panel for troubleshooting
- Automatic initial state loading on connection
- Optimized rendering for high-frequency updates

**Note**: The React dashboard is a demonstration tool to showcase SSE capabilities. The core SSE implementation is server-side and can be used with any client.

### Stopping the Server

Press `Ctrl+C` to initiate graceful shutdown. The server will:
1. Stop accepting new connections
2. Stop event retention
3. Wait up to 1 minute for existing connections to close
4. Shut down gracefully

## API Examples

### SSE Endpoint

- `GET /events/watch` - SSE endpoint for real-time events
  - Optional header: `Last-Event-ID` - Resume from a specific event ID

### Sample Domain Endpoints (Metrics)

For detailed API documentation and examples, see the `docs/api/` directory:

- `docs/api/metrics_api_docs.http` - Metrics API examples
- `docs/api/metric_readings_api_docs.http` - Metric readings API examples

These endpoints demonstrate how domain actions trigger SSE events:
- `POST /metrics` - Create a metric (triggers `metric_created` event)
- `GET /metrics` - Get all metrics (with optional `?with_readings=true` parameter)
- `GET /metrics/:id` - Get metric by ID
- `POST /metrics/readings` - Create a reading (triggers `metric_reading_created` event)

A React dashboard is available in `cmd/client/` to visually demonstrate SSE in action.

## SSE Implementation

**Reference**: For a comprehensive guide on Server-Sent Events, see the [MDN documentation on Using server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events).

### How It Works

1. **Client Connection**: Clients connect to `/events/watch` endpoint
2. **Event Broadcasting**: When domain events occur, they're broadcast to all connected clients
3. **Event Storage**: Events are stored in the event store for replay
4. **Reconnection Support**: Clients can use `Last-Event-ID` header to receive missed events

### Event Flow

```
Domain Action → Use Case → SSE Hub.Broadcast → Event Store
                                          ↓
                                    All Connected Clients
```

### Key Components

- **SSE Hub** (`pkg/sse/sse_hub.go`): Manages client connections and broadcasting
- **Event Store** (`pkg/sse/event_store.go`): Interface for event storage and replay
- **SSE Client** (`pkg/sse/client.go`): Internal client representation

### Event Replay

Clients can reconnect using the `Last-Event-ID` header to receive events that occurred while disconnected. The event store maintains events for a configurable TTL (default: 1 minute).

### Connection Management

- Maximum 10,000 concurrent clients (configurable)
- Oldest client disconnected when limit reached
- Slow clients automatically dropped if they can't keep up with the event stream
- Graceful connection cleanup on disconnect

See `cmd/client/` for a complete React dashboard implementation demonstrating SSE connectivity.

## Architecture Details

### SSE Hub (`pkg/sse/sse_hub.go`)

The core component that manages all SSE functionality. Implemented as a singleton pattern for global access:

- **Singleton Pattern**: Initialized once with `InitializeSSEHub()` and accessed via `GetSSEHub()`
- **Client Registration**: Registers new SSE clients via `Register` channel
- **Client Unregistration**: Removes disconnected clients via `Unregister` channel
- **Event Broadcasting**: Receives events via `Broadcast` channel and sends to all clients
- **Client Limit**: Maximum concurrent clients (default: 10,000) - oldest disconnected when limit reached
- **Slow Client Handling**: Non-blocking sends - drops clients if their channel is full
- **Thread-Safe**: Uses channels and `sync.Once` for safe concurrent operations

**Initialization**: The SSE Hub is initialized during application startup in `main.go` with the event store and max clients configuration. Controllers and use cases access it via `sse.GetSSEHub()`.

### Event Store (`pkg/sse/event_store.go`)

Interface for event storage and replay:

- **StoreEvent**: Stores events for later replay
- **GetEventsAfterID**: Retrieves events after a given ID for reconnection support
- **TTL-based Retention**: Events automatically expire after TTL (default: 1 minute)
- **Thread-Safe**: In-memory implementation uses mutexes for safe concurrent access

### Sample Domain (Metrics)

The metrics domain is provided as a demonstration of how to integrate SSE with domain logic:

- **Entities**: `Metric` and `MetricReading` (sample domain entities)
- **Use Cases**: Trigger SSE events when domain actions occur
- **Controllers**: HTTP endpoints that trigger domain actions, which in turn broadcast SSE events

The SSE infrastructure is completely independent of the metrics domain and can be used with any domain.

## Configuration

Default configuration (in `cmd/server/main.go`):
- **Port**: `8089`
- **Max SSE Clients**: `10,000`
- **Event TTL**: `1 minute`
- **Graceful Shutdown Timeout**: `1 minute`

**SSE Hub Initialization**: The SSE Hub singleton is initialized during application startup via `sse.InitializeSSEHub(eventStore, maxClients)`. It must be initialized before any components attempt to access it via `sse.GetSSEHub()`.

To modify these values, edit the constants and variables in `main.go`.

## Demonstration

### Visual Demo

A React dashboard is included to demonstrate the SSE system in action. The dashboard shows:

- **Real-time Updates**: Metrics and readings appear instantly via SSE without polling
- **Initial State Loading**: New connections automatically load existing data
- **High-Frequency Support**: Optimized rendering handles updates faster than 500ms
- **Adaptive Animations**: Chart animations adjust based on metric input frequency
- **Connection Monitoring**: Visual indicators for connection status and health

**Video/GIF Demo**:

https://github.com/user-attachments/assets/da8d5a7c-87ba-40ab-aab5-d23c8c35396a

**Note**: The React dashboard was created by Auto (Cursor AI agent) as a demonstration tool to showcase SSE capabilities. The dashboard serves as a visual demonstration tool. The SSE implementation itself is client-agnostic and can be integrated with any front-end technology.

## Development

### Building

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

### Testing

The project follows Go best practices with interface-based design for testability. Use the provided HTTP files in `docs/api/` for API testing with REST Client extensions or tools like Postman, cURL, or HTTPie.

## Learning Journey

This project was developed as a learning exercise, starting from basic concepts and progressing to more advanced implementations. The development conversation with ChatGPT documents part of the learning process, from initial simple questions to more complex architectural decisions:

**[Development Conversation with ChatGPT](https://chatgpt.com/share/6950a4be-77e4-8003-9c73-a23b22038ae6)**

This conversation may be helpful for understanding:
- The learning progression from beginner to a little more advanced topics
- Design decisions and their rationale
- Problem-solving approaches throughout development
- Final outcomes and implementation details

Naturally, I didn't vibe-coded the backend lol. I also did some research on regular sources (e.g. StackOverflow), specially for GoLang-specific gotchas.

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
