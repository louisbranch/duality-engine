# Getting Started

## Prerequisites

- Go 1.21+

## Run locally

Start the gRPC server and MCP bridge together:

```sh
make run
```

This runs the gRPC server on `localhost:8080`, waits for it to accept
connections, and then starts the MCP server on stdio.

## Run services individually

Start the gRPC server:

```sh
go run ./cmd/server
```

Start the MCP server (requires the gRPC server running):

```sh
go run ./cmd/mcp
```
