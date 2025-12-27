# Go SSE Sample with DDD

A Go-based Server-Sent Events (SSE) sample application demonstrating Domain-Driven Design (DDD) architecture. This project provides a RESTful API for managing metrics and metric readings, with real-time event broadcasting through SSE.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [SSE Events](#sse-events)
- [Architecture Details](#architecture-details)
- [Configuration](#configuration)
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

A Node.js example client is provided to demonstrate SSE connectivity. Make sure the server is running first.

Using Make:
```bash
make run-client
```

Or directly with Node.js:
```bash
node cmd/client/client.mjs
```

The client will:
- Connect to the SSE endpoint
- Listen for `metric_created` and `metric_reading_created` events
- Log events to the console as they are received

### Stopping the Server

Press `Ctrl+C` to initiate graceful shutdown. The server will:
1. Stop accepting new connections
2. Stop event retention
3. Wait up to 1 minute for existing connections to close
4. Shut down gracefully

## API Documentation

### Base URL

```
http://localhost:8089
```

### Metrics API

#### Create Metric

Creates a new metric with a unique identifier.

**Endpoint:** `POST /metrics`

**Request Body:**
```json
{
  "name": "CPU Temperature"
}
```

**Response:** `201 Created`
```json
{
  "id": "018f1234-5678-7890-abcd-ef1234567890",
  "name": "CPU Temperature"
}
```

**Validation:**
- `name` is required and cannot be empty

#### Get Metric by ID

Retrieves a metric by its unique identifier.

**Endpoint:** `GET /metrics/:id`

**Path Parameters:**
- `id`: UUID v7 of the metric

**Response:** `200 OK`
```json
{
  "id": "018f1234-5678-7890-abcd-ef1234567890",
  "name": "CPU Temperature"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `500 Internal Server Error`: Metric not found or server error

### Metric Readings API

#### Create Metric Reading

Creates a new reading for an existing metric.

**Endpoint:** `POST /metrics/readings`

**Request Body:**
```json
{
  "metric_id": "018f1234-5678-7890-abcd-ef1234567890",
  "value": 45.5,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Fields:**
- `metric_id` (required): UUID v7 of the metric
- `value` (required): Numeric value, must be greater than 0
- `timestamp` (optional): RFC3339 formatted timestamp. If omitted, current UTC time is used

**Response:** `201 Created`
```json
{
  "id": "018f1234-5678-7890-abcd-ef1234567891",
  "metric_id": "018f1234-5678-7890-abcd-ef1234567890",
  "value": 45.5,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Validation:**
- `metric_id` must be a valid UUID v7
- `value` must be greater than 0
- `timestamp` must be in RFC3339 format if provided
- The metric must exist

**Error Responses:**
- `400 Bad Request`: Invalid input (UUID, timestamp format, or value validation)
- `500 Internal Server Error`: Metric not found or server error

### Events API (SSE)

#### Watch Events

Establishes a Server-Sent Events connection to receive real-time notifications.

**Endpoint:** `GET /events/watch`

**Headers:**
- `Last-Event-ID` (optional): Event ID to resume from. Events after this ID will be sent immediately upon connection.

**Response:** `200 OK` (streaming)

**Content-Type:** `text/event-stream`

**Connection Behavior:**
- Connection remains open until client disconnects
- Automatic reconnection support via `Last-Event-ID` header
- Connection messages are sent upon connect/disconnect

**Example Event:**
```json
{
  "id": "018f1234-5678-7890-abcd-ef1234567892",
  "type": "metric_created",
  "data": {
    "id": "018f1234-5678-7890-abcd-ef1234567890",
    "name": "CPU Temperature"
  }
}
```

**Client Example (Browser JavaScript):**
```javascript
const eventSource = new EventSource('http://localhost:8089/events/watch');

eventSource.addEventListener('metric_created', (event) => {
  const data = JSON.parse(event.data);
  console.log('New metric created:', data);
});

eventSource.addEventListener('metric_reading_created', (event) => {
  const data = JSON.parse(event.data);
  console.log('New reading created:', data);
});

// Handle reconnection with Last-Event-ID
eventSource.onerror = (error) => {
  console.error('SSE error:', error);
  // EventSource automatically handles reconnection
};
```

**Client Example (Node.js):**

A complete Node.js client example is available in `cmd/client/client.mjs`:

```javascript
import { EventSource } from "eventsource";

const es = new EventSource("http://localhost:8089/events/watch");

es.onmessage = (event) => {
  switch (event.data) {
    case "connected": {
      console.log("SSE connected");
      break;
    }
    case "disconnected": {
      console.log("SSE disconnected");
      es.close();
      break;
    }
    default: {
      console.log("Unexpected SSE message:", event.data);
      break;
    }
  }
};

es.addEventListener("metric_created", (event) => {
  const data = JSON.parse(event.data);
  console.log("metric created:", data);
});

es.addEventListener("metric_reading_created", (event) => {
  const data = JSON.parse(event.data);
  console.log("metric reading created:", data);
});

es.onerror = (err) => {
  console.error("SSE error:", err);
};
```

Run it with:
```bash
make run-client
# or
node cmd/client/client.mjs
```

**Client Example (cURL):**
```bash
curl -N -H "Accept: text/event-stream" http://localhost:8089/events/watch
```

## SSE Events

### Event Types

The application broadcasts two types of events:

1. **`metric_created`**: Emitted when a new metric is created
   ```json
   {
     "id": "018f1234-5678-7890-abcd-ef1234567892",
     "type": "metric_created",
     "data": {
       "id": "018f1234-5678-7890-abcd-ef1234567890",
       "name": "CPU Temperature"
     }
   }
   ```

2. **`metric_reading_created`**: Emitted when a new metric reading is created
   ```json
   {
     "id": "018f1234-5678-7890-abcd-ef1234567893",
     "type": "metric_reading_created",
     "data": {
       "id": "018f1234-5678-7890-abcd-ef1234567891",
       "metric_id": "018f1234-5678-7890-abcd-ef1234567890",
       "value": 45.5,
       "timestamp": "2024-01-15T10:30:00Z"
     }
   }
   ```

### Event Replay

When a client connects with the `Last-Event-ID` header, the server will:
1. Retrieve all events stored after the specified event ID
2. Send them immediately upon connection
3. Continue streaming new events as they occur

This enables clients to recover from disconnections without missing events.

### Event Retention

Events are stored in memory with a configurable TTL (Time To Live). By default:
- TTL: 1 minute
- Retention runs every 30 seconds (TTL/2)
- Events older than TTL are automatically removed

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

### Server Configuration

Default configuration in `cmd/server/main.go`:

- **Port**: `8089`
- **Max SSE Clients**: `10,000`
- **Event TTL**: `1 minute`
- **Graceful Shutdown Timeout**: `1 minute`

To modify these values, edit the constants and variables in `main.go`.

### Environment Variables

Currently, the application uses hardcoded configuration. To make it configurable via environment variables, you can:

1. Add environment variable parsing
2. Use a configuration library (e.g., `viper`, `envconfig`)
3. Provide default values

## Development

### Project Dependencies

**Server dependencies:**
- `github.com/gin-gonic/gin`: Web framework
- `github.com/google/uuid`: UUID generation

**Client dependencies:**
- `eventsource`: Node.js Server-Sent Events client library

### Code Organization

The project follows Go best practices:
- Package naming conventions
- Interface-based design for testability
- Clear separation of concerns
- Thread-safe implementations

### Testing

To add tests:
1. Create test files with `_test.go` suffix
2. Use Go's built-in testing package
3. Mock repositories for unit testing
4. Integration tests for API endpoints

### Building

Build the application:
```bash
go build -o bin/server cmd/server/main.go
```

Run the binary:
```bash
./bin/server
```

### API Testing

Use the provided HTTP files in `docs/api/` with REST Client extensions (VS Code, IntelliJ) or tools like:
- Postman
- cURL
- HTTPie

Example using the provided HTTP files:
```http
@baseUrl = http://localhost:8089

POST {{baseUrl}}/metrics
Content-Type: application/json

{
    "name": "CPU Temperature"
}
```

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
