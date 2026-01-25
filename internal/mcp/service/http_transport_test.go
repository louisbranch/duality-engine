package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
)

func TestHTTPTransport_Connect(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	ctx := context.Background()

	conn, err := transport.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	if conn == nil {
		t.Fatal("Connect() returned nil connection")
	}

	sessionID := conn.SessionID()
	if sessionID == "" {
		t.Error("SessionID() returned empty string")
	}

	// Test that connection can be closed
	if err := conn.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestHTTPTransport_handleHealth(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	
	req := httptest.NewRequest(http.MethodGet, "/mcp/health", nil)
	w := httptest.NewRecorder()
	
	transport.handleHealth(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("handleHealth() status = %d, want %d", w.Code, http.StatusOK)
	}
	
	if w.Body.String() != "OK" {
		t.Errorf("handleHealth() body = %q, want %q", w.Body.String(), "OK")
	}
}

func TestHTTPTransport_handleMessages_InvalidMethod(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	
	req := httptest.NewRequest(http.MethodGet, "/mcp/messages", nil)
	w := httptest.NewRecorder()
	
	transport.handleMessages(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("handleMessages() status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHTTPTransport_handleMessages_InvalidJSON(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	
	req := httptest.NewRequest(http.MethodPost, "/mcp/messages", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	transport.handleMessages(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("handleMessages() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHTTPTransport_handleMessages_NewSession(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	
	// Create a simple JSON-RPC request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params":  map[string]interface{}{},
	}
	body, _ := json.Marshal(request)
	
	req := httptest.NewRequest(http.MethodPost, "/mcp/messages", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	// This will create a new session but won't get a response without a running MCP server
	// So we expect it to timeout or handle gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	
	transport.handleMessages(w, req)
	
	// Should have created a session (check header)
	sessionID := w.Header().Get("X-MCP-Session-ID")
	if sessionID == "" {
		t.Error("handleMessages() should set X-MCP-Session-ID header for new sessions")
	}
}

func TestHTTPTransport_handleSSE_InvalidMethod(t *testing.T) {
	transport := NewHTTPTransport("localhost:8081")
	
	req := httptest.NewRequest(http.MethodPost, "/mcp/sse", nil)
	w := httptest.NewRecorder()
	
	transport.handleSSE(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("handleSSE() status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHTTPConnection_ReadWrite(t *testing.T) {
	conn := &httpConnection{
		sessionID: "test_session",
		reqChan:   make(chan jsonrpc.Message, 1),
		respChan:  make(chan jsonrpc.Message, 1),
		closed:    make(chan struct{}),
	}
	
	ctx := context.Background()
	
	// Test Write
	request := &jsonrpc.Request{
		Method: "test",
		ID:     jsonrpc.ID{},
	}
	if err := conn.Write(ctx, request); err != nil {
		t.Errorf("Write() error = %v", err)
	}
	
	// Test Read (should read what we wrote, but Read reads from reqChan, not respChan)
	// So we need to write to reqChan directly
	conn.reqChan <- request
	
	msg, err := conn.Read(ctx)
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
	if msg == nil {
		t.Error("Read() returned nil message")
	}
}

func TestHTTPConnection_Close(t *testing.T) {
	conn := &httpConnection{
		sessionID: "test_session",
		reqChan:   make(chan jsonrpc.Message, 1),
		respChan:  make(chan jsonrpc.Message, 1),
		closed:    make(chan struct{}),
	}
	
	// Close should not error
	if err := conn.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	
	// Close again should also not error
	if err := conn.Close(); err != nil {
		t.Errorf("Close() second call error = %v", err)
	}
	
	// Write after close should error
	ctx := context.Background()
	request := &jsonrpc.Request{
		Method: "test",
		ID:     jsonrpc.ID{},
	}
	if err := conn.Write(ctx, request); err == nil {
		t.Error("Write() after Close() should return error")
	}
}

func TestHTTPConnection_SessionID(t *testing.T) {
	conn := &httpConnection{
		sessionID: "test_session_123",
		reqChan:   make(chan jsonrpc.Message, 1),
		respChan:  make(chan jsonrpc.Message, 1),
		closed:    make(chan struct{}),
	}
	
	if got := conn.SessionID(); got != "test_session_123" {
		t.Errorf("SessionID() = %q, want %q", got, "test_session_123")
	}
}

func TestSessionTransport_Connect(t *testing.T) {
	conn := &httpConnection{
		sessionID: "test_session",
		reqChan:   make(chan jsonrpc.Message, 1),
		respChan:  make(chan jsonrpc.Message, 1),
		closed:    make(chan struct{}),
	}
	
	transport := &sessionTransport{conn: conn}
	ctx := context.Background()
	
	returnedConn, err := transport.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	
	if returnedConn != conn {
		t.Error("Connect() returned different connection")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1 := generateSessionID()
	id2 := generateSessionID()
	
	if id1 == id2 {
		t.Error("generateSessionID() should generate unique IDs")
	}
	
	if !strings.HasPrefix(id1, "session_") {
		t.Errorf("generateSessionID() = %q, should start with 'session_'", id1)
	}
}
