# Gophormula Architecture Review

## Project Vision

Gophormula is a Go application for ingesting, parsing, and displaying Formula 1 race telemetry and timing data. It bridges F1's live timing infrastructure with consumption layers (dashboards, metrics, analysis). The vision is a comprehensive F1 telemetry platform supporting:

- **Live streaming** via SignalR from `livetiming.formula1.com`
- **Historical replay** of archived race sessions
- **Dashboards** (SSE-driven web UI, TUI, Grafana/Prometheus)
- **Data archival** for offline analysis

---

## Package Structure

```
gophormula/
├── cmd/
│   ├── live/      # Live SignalR consumer (log-only)
│   ├── replay/    # Historical replay runner (log-only)
│   ├── historic/  # Data archival downloader
│   └── dash/      # Web dashboard server
└── pkg/
    ├── signalr/    # SignalR protocol client
    ├── livetiming/ # F1 data parsing, type definitions, circuit map fetch
    ├── replay/     # File-based pub-sub replay engine with seek support
    └── dash/       # SSE-driven web dashboard (Hub, standings, status, SVG track map)
```

### Package Responsibilities

| Package | Responsibility |
|---------|---------------|
| `pkg/signalr` | SignalR protocol (negotiation, handshake, WebSocket transport) |
| `pkg/livetiming` | F1 message types, decompression, reflection-based parser, Multiviewer circuit map client |
| `pkg/replay` | File-based pub-sub replay with seek/fast-forward and subscriber channels |
| `pkg/dash` | SSE hub, HTTP handlers, standings accumulator, status accumulator, SVG track renderer |

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
              │  - SSE Web Dashboard    │
              │  - Metrics export       │
              │  - Analysis tools       │
              └─────────────────────────┘
```

---

## Architectural Strengths

### 1. Clean Separation of Concerns
Transport, protocol, parsing, and application logic are cleanly separated. Each package has a single, focused responsibility. The SignalR client knows nothing about F1 data; the parser knows nothing about network transport.

### 2. Intentionally Minimal Dependencies
Only two external dependencies (`gorilla/websocket`, `datastar-go`). Heavy use of the standard library reduces maintenance burden and avoids dependency hell — a solid choice for a long-lived project.

### 3. Extensible Abstractions
- **Transport interface** makes swapping WebSocket for SSE or long-polling straightforward
- **Options pattern** (`WithURL`, `WithVersion`) in the SignalR client is clean and extensible
- **Reflection-based message registry** allows registering new message types without changing dispatch logic

### 4. Streaming / Memory-Conscious Design
Bufio scanners for line-by-line file reading, streaming JSON parsing from HTTP — the codebase avoids loading entire race sessions into memory, which matters given race files can be 20-50MB per session.

### 5. Concurrent Download Model
`cmd/historic` uses `sync.WaitGroup` for concurrent feed downloads with clean goroutine coordination.

### 6. Incremental State Accumulation
The dashboard accumulates `TimingData`, `DriverList`, and status bar values server-side so that any newly connected client receives the full current state, and so that seek/fast-forward correctly seeds standings before real-time playback begins.

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
Mix of `log` and `log/slog` across packages. No configurable log levels. Hard to debug in production.

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

The project has a strong foundation and a clear, realistic vision. The core loop — download → parse → replay/display — is architecturally sound and fully working end-to-end. The hardest part (reverse-engineering the F1 SignalR protocol and data formats, including `.z` compressed streams) is done. The SSE dashboard delivers a real-time race replay UI with live standings, circuit track map (via Multiviewer API), weather, lap count, and session status.

What remains is mostly hardening and expansion:

- Filling in incomplete implementations (live reconnection)
- Extending the consumption layer (Prometheus metrics, TUI)
- Hardening for long-running operation (context propagation, graceful shutdown, observability)

The minimal-dependency philosophy and clean package boundaries are commendable and should be preserved as the project grows.