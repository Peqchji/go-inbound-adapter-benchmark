# Go Inbound Adapter Benchmark

## Overview

This project is a benchmark suite designed to compare the performance and characteristics of different inbound adapters in a Go application. It aims to implement and measure:
- gRPC
- GraphQL
- HTTP

The domain logic centers around a simple Wallet service.

## Project Status

**Work In Progress**: This repository currently contains the skeletal structure and some basic domain/persistence logic. The server implementations (gRPC, GraphQL, HTTP) and the actual benchmark runner are yet to be populated.

## Project Structure

The project follows a standard Go project layout:

- `cmd/`: Application entry points.
- `internal/`: Private application and library code.
  - `client/`: External interface implementations (e.g., Database).
  - `domain/`: Business logic (Wallet entity, repository interfaces).
  - `server/`: Inbound adapter implementations (Placeholders for gql, grpc, http).
- `pkg/`: Library code that's ok to use by external applications (e.g., shared interfaces/types).

## Getting Started

*(Instructions will be added as the project matures)*
