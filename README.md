# Go Inbound Adapter Benchmark

## Overview

This project benchmarks the performance of three different inbound adapter protocols for a Go application:
1. **REST (HTTP)**: Implemented using [Echo v4](https://echo.labstack.com/).
2. **GraphQL**: Implemented using [gqlgen](https://github.com/99designs/gqlgen) wrapped in Echo v4.
3. **gRPC**: Implemented using native [gRPC-Go](https://github.com/grpc/grpc-go).

All three adapters share the same **Domain Layer** (In-Memory Wallet Service) to ensure the benchmark measures protocol/transport overhead, not business logic.

## Architecture

- **Domain**: Pure business logic (`internal/domain/wallet`).
- **Adapter**: Protocol-specific handlers (REST, GQL, gRPC) in `internal/adapter` and `cmd`.
- **Fairness Strategy**: We wrapped both REST and GraphQL servers in the **Echo** framework with identical middleware (`Logger`, `Recover`) to simulate a realistic production HTTP environment. gRPC runs on its native server stack.

## Prerequisites

- **Go** (1.23+)
- **k6** (for running benchmarks)
- **PowerShell** (for the start script)

## How to Run

### 1. Start All Servers
We provide a helper script to start all three servers on different ports:
- REST: `:8080`
- GraphQL: `:8081`
- gRPC: `:8082`

```powershell
docker compose up --build
```

### 2. Run k6 Benchmark
The benchmark suite tests a full "Create Wallet -> Get Wallet" flow for each protocol.

```powershell
k6 run k6/suite.js --summary-export=summary.json
```

## Benchmark Results (Typical)

| Adapter | Tech Stack | Avg Response Time | P95 Latency |
|---------|------------|-------------------|-------------|
| **gRPC** | Native / HTTP/2 | ~1.2 ms | ~2.5 ms |
| **REST** | Echo v4 | ~5.1 ms | ~11.3 ms |
| **GraphQL** | Echo v4 + gqlgen | ~5.5 ms | ~12.0 ms |

## Project Structure

- `cmd/`: Entry points for each server (`gqlserver`, `grpcserver`, `httpserver`).
- `internal/`:
  - `adapter/`: Handlers for REST, In-Memory Repo.
  - `domain/`: Wallet entities and Service logic.
- `k6/`: Load testing scripts.

