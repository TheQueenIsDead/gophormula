Can # Gophormula Architecture Review

## Project Vision

Gophormula is a Go application for ingesting, parsing, and displaying Formula 1 race telemetry and timing data. It bridges F1's live timing infrastructure with consumption layers (dashboards, metrics, analysis). The vision is a comprehensive F1 telemetry platform supporting:

- **Live streaming** via SignalR from `livetiming.formula1.com`
- **Historical replay** of archived race sessions
- **Dashboards** (Grafana/Prometheus, SSE-driven web UI, TUI)
- **Data archival** for offline analysis

---

## Package Structure

```
gophormula/
├── cmd/
│   ├── live/      # Live SignalR consumer
│   ├── replay/    # Historical replay runner
│   ├── historic/  # Data archival downloader
│   └── dash/      # Web dashboard server
└── pkg/
    ├── signalr/   # SignalR protocol client
    ├── livetiming/ # F1 data parsing & type definitions
    ├── replay/    # Pub-sub replay engine
    └── dash/      # HTTP handler (stub)
```

### Package Responsibilities

| Package | Responsibility |
|---------|---------------|
| `pkg/signalr` | SignalR protocol (negotiation, handshake, WebSocket transport) |
| `pkg/livetiming` | F1 message types, decompression, reflection-based parser |
| `pkg/replay` | File-based pub-sub replay with subscriber channels |
| `pkg/dash` | Web server (early stub) |

---

## Data Flow

```
┌─────────────────────────────────────────────────────┐
│                   Data Sources                       │
├──────────────────────┬──────────────────────────────┤
│  Live Stream         │  Historical Archives          │
│ (livetiming.f1.com)  │ (livetiming.f1.com/static)  │
└──────────┬───────────┴──────────────┬───────────────┘
           │                          │
    ┌──────▼──────────┐    ┌──────────▼──────────┐
    │  SignalR Client  │    │  HTTP Downloader    │
    │  (live cmd)      │    │  (historic cmd)     │
    └──────┬───────────┘    └──────────┬──────────┘
           │                           │
           │                 ┌─────────▼──────────┐
           │                 │  Local JSON Files  │
           │                 └─────────┬──────────┘
           │                           │
           └─────────────┬─────────────┘
                         │
              ┌──────────▼──────────────┐
              │   Livetiming Parser     │
              │  - Topic routing        │
              │  - Flate decompression  │
              │  - JSON unmarshalling   │
              │  - Reflection dispatch  │
              └──────────┬──────────────┘
                         │
              ┌──────────▼──────────────┐
              │  Consumer Applications  │
              │  - Dashboards           │
              │  - Metrics export       │
              │  - Analysis tools       │
              └─────────────────────────┘
```

---

## Architectural Strengths

### 1. Clean Separation of Concerns
Transport, protocol, parsing, and application logic are cleanly separated. Each package has a single, focused responsibility. The SignalR client knows nothing about F1 data; the parser knows nothing about network transport.

### 2. Intentionally Minimal Dependencies
Only one external dependency (`gorilla/websocket`). Heavy use of the standard library (`encoding/json`, `compress/flate`, `encoding/base64`, `log/slog`) reduces maintenance burden and avoids dependency hell — a solid choice for a long-lived project.

### 3. Extensible Abstractions
- **Transport interface** makes swapping WebSocket for SSE or long-polling straightforward
- **Options pattern** (`WithURL`, `WithVersion`) in the SignalR client is clean and extensible
- **Reflection-based message registry** allows registering new message types without changing dispatch logic

### 4. Streaming / Memory-Conscious Design
Bufio scanners for line-by-line file reading, streaming JSON parsing from HTTP — the codebase avoids loading entire race sessions into memory, which matters given race files can be 20-50MB per session.

### 5. Concurrent Download Model
`cmd/historic` uses `sync.WaitGroup` for concurrent feed downloads with clean goroutine coordination.

---

## Architectural Weaknesses

### Critical (Build/Runtime Correctness)

**1. Stubbed reconnection logic** (`pkg/signalr/client.go`)
`reconnect()`, `abort()`, and `ping()` all `panic("not implemented")`. A live F1 session runs for hours; without reconnection the client is not production-usable.

### Design Gaps

**2. Pointer-to-channel pattern** (`pkg/replay/replay.go`)
`subscribers []*chan any` takes the address of a channel, which is unusual and unnecessary — channels are already reference types. This adds confusion without benefit.

### Observability

**3. Inconsistent logging**
Mix of `log`, `log/slog`, and silent errors across packages. No configurable log levels. Hard to debug in production.

**4. No metrics**
Prometheus integration is in the roadmap but absent. No message counters, error rates, or latency tracking anywhere.

**5. No context/timeout support**
HTTP operations in `cmd/historic` and WebSocket operations in `pkg/signalr` have no context propagation. Hangs on slow/dead connections cannot be cancelled.

### Test Coverage

**6. Very limited tests**
Only `pkg/livetiming/parser_test.go` exists (decompression tests). No tests for SignalR protocol, replay timing, or message type registration. Integration tests against recorded sessions would provide high confidence.

---

## Recommendations

### Priority 1 — Correctness
1. Implement `reconnect()` with exponential backoff (not panic)
2. Add context/timeout to all HTTP and WebSocket operations

### Priority 2 — Architecture
3. Remove pointer-to-channel anti-pattern — use `[]chan any` directly

### Priority 3 — Observability
4. Settle on `log/slog` throughout, with configurable level via env/flag
5. Add Prometheus metrics: messages parsed, errors, connection uptime, replay lag
6. Add health/readiness endpoints to `cmd/dash`

### Priority 4 — Quality
7. Add table-driven tests for parser with real recorded SignalR payloads
8. Add graceful shutdown via `os.Signal` + context cancellation across all cmds

---

## Vision Assessment

The project has a strong foundation and a clear, realistic vision. The core loop — download → parse → replay/display — is architecturally sound. The hardest part (reverse-engineering the F1 SignalR protocol and data formats) is largely done. What remains is mostly completeness and polish:

- Filling in incomplete implementations (parsing, reconnection, replay timing)
- Building out the consumption layer (dashboard, metrics)
- Hardening for long-running operation (reconnection, error handling, observability)

The minimal-dependency philosophy and clean package boundaries are commendable and should be preserved as the project grows. Resist the temptation to add a large framework — the stdlib-first approach is a genuine strength here.