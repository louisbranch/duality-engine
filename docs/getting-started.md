# Getting Started

## Prerequisites

- Go 1.25.6
- protoc (required until binaries are published)
- BoltDB (embedded; the server creates `data/duality.db` by default)
- Make (for `make run`)

## Run locally

Start the gRPC server and MCP bridge together:

```sh
make run
```

This runs the gRPC server on `localhost:8080` and the MCP server on stdio.
The MCP server will wait for the gRPC server to be healthy before accepting requests.

## Run services individually

Start the gRPC server:

```sh
go run ./cmd/server
```

Start the MCP server after the gRPC server starts. See
[MCP tools and resources](mcp.md) for the MCP run command.
