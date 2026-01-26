#!/bin/bash
# Test script for MCP HTTP transport initialization sequence
# This script simulates the proper MCP protocol flow:
# 1. initialize request -> get session ID
# 2. initialized notification -> complete initialization
# 3. tools/list request -> verify server is ready

set -e

MCP_URL="${MCP_URL:-http://localhost:3001/mcp}"

echo "=== MCP HTTP Transport Test ==="
echo "Testing endpoint: $MCP_URL"
echo ""

# Step 1: Send initialize request
echo "Step 1: Sending initialize request..."
INIT_RESPONSE=$(curl -sS -D /tmp/mcp-headers.txt -X POST "$MCP_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "clientInfo": {
        "name": "test-client",
        "version": "0.1.0"
      }
    }
  }')

echo "Initialize response: $INIT_RESPONSE"
echo ""

# Extract session ID from response headers
SESSION_ID=$(grep -i "X-MCP-Session-ID" /tmp/mcp-headers.txt | cut -d' ' -f2 | tr -d '\r' | tr -d '\n')

if [ -z "$SESSION_ID" ]; then
  echo "ERROR: No session ID found in response headers"
  echo "Response headers:"
  cat /tmp/mcp-headers.txt
  exit 1
fi

echo "Session ID: $SESSION_ID"
echo ""

# Step 2: Send initialized notification
echo "Step 2: Sending initialized notification..."
INITIALIZED_RESPONSE=$(curl -sS -w "\nHTTP Status: %{http_code}\n" -X POST "$MCP_URL" \
  -H "Content-Type: application/json" \
  -H "X-MCP-Session-ID: $SESSION_ID" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialized",
    "params": {}
  }')

echo "Initialized response: $INITIALIZED_RESPONSE"
echo ""

# Step 3: Send tools/list request to verify server is ready
echo "Step 3: Sending tools/list request..."
TOOLS_LIST_RESPONSE=$(curl -sS -X POST "$MCP_URL" \
  -H "Content-Type: application/json" \
  -H "X-MCP-Session-ID: $SESSION_ID" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list",
    "params": {}
  }')

echo "Tools list response: $TOOLS_LIST_RESPONSE"
echo ""

# Verify the response
if echo "$TOOLS_LIST_RESPONSE" | grep -q '"result"'; then
  echo "✓ Success: Server is initialized and responding to requests"
  exit 0
else
  echo "✗ Error: Server did not respond correctly"
  echo "Response: $TOOLS_LIST_RESPONSE"
  exit 1
fi
