// Package mcp tests the MCP server wiring.
package mcp

import (
	"context"
	"errors"
	"testing"

	pb "github.com/louisbranch/duality-protocol/api/gen/go/duality/v1"
	"github.com/mark3labs/mcp-go/mcp"
	"google.golang.org/grpc"
)

// fakeDiceRollClient implements DiceRollServiceClient for tests.
type fakeDiceRollClient struct {
	response            *pb.ActionRollResponse
	rollDiceResponse    *pb.RollDiceResponse
	err                 error
	rollDiceErr         error
	lastRequest         *pb.ActionRollRequest
	lastRollDiceRequest *pb.RollDiceRequest
}

// ActionRoll records the request and returns the configured response.
func (f *fakeDiceRollClient) ActionRoll(ctx context.Context, req *pb.ActionRollRequest, opts ...grpc.CallOption) (*pb.ActionRollResponse, error) {
	f.lastRequest = req
	return f.response, f.err
}

// RollDice records the request and returns the configured response.
func (f *fakeDiceRollClient) RollDice(ctx context.Context, req *pb.RollDiceRequest, opts ...grpc.CallOption) (*pb.RollDiceResponse, error) {
	f.lastRollDiceRequest = req
	return f.rollDiceResponse, f.rollDiceErr
}

// newCallToolRequest builds a tool call request with arguments.
func newCallToolRequest(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "duality_action_roll",
			Arguments: args,
		},
	}
}

// newRollDiceCallToolRequest builds a roll dice tool call request with arguments.
func newRollDiceCallToolRequest(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "roll_dice",
			Arguments: args,
		},
	}
}

// TestGRPCAddressPrefersEnv ensures env configuration overrides defaults.
func TestGRPCAddressPrefersEnv(t *testing.T) {
	t.Setenv("DUALITY_GRPC_ADDR", "env:123")
	if got := grpcAddress("fallback"); got != "env:123" {
		t.Fatalf("expected env address, got %q", got)
	}
}

// TestGRPCAddressFallback ensures the fallback address is used when env is empty.
func TestGRPCAddressFallback(t *testing.T) {
	t.Setenv("DUALITY_GRPC_ADDR", "")
	if got := grpcAddress("fallback"); got != "fallback" {
		t.Fatalf("expected fallback address, got %q", got)
	}
}

// TestServeRequiresConfiguredServer ensures Serve returns an error when unconfigured.
func TestServeRequiresConfiguredServer(t *testing.T) {
	tests := []struct {
		name   string
		server *Server
	}{
		{name: "nil server", server: nil},
		{name: "missing mcp server", server: &Server{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.server.Serve(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

// TestNewConfiguresServer ensures New returns a configured server.
func TestNewConfiguresServer(t *testing.T) {
	server, err := New("localhost:8080")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if server == nil || server.mcpServer == nil {
		t.Fatal("expected configured server")
	}
}

// TestActionRollHandlerRejectsNegativeDifficulty ensures invalid difficulty returns an error result.
func TestActionRollHandlerRejectsNegativeDifficulty(t *testing.T) {
	client := &fakeDiceRollClient{}
	handler := actionRollHandler(client)

	result, err := handler(context.Background(), newCallToolRequest(map[string]any{
		"modifier":   1,
		"difficulty": -1,
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
	if client.lastRequest != nil {
		t.Fatal("expected no gRPC call on invalid input")
	}
}

// TestActionRollHandlerReturnsClientError ensures gRPC errors are returned as tool errors.
func TestActionRollHandlerReturnsClientError(t *testing.T) {
	client := &fakeDiceRollClient{err: errors.New("boom")}
	handler := actionRollHandler(client)

	result, err := handler(context.Background(), newCallToolRequest(map[string]any{
		"modifier": 2,
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
}

// TestActionRollHandlerHandlesMissingDice ensures missing dice results in an error result.
func TestActionRollHandlerHandlesMissingDice(t *testing.T) {
	client := &fakeDiceRollClient{
		response: &pb.ActionRollResponse{
			Outcome: pb.Outcome_OUTCOME_UNSPECIFIED,
		},
	}
	handler := actionRollHandler(client)

	result, err := handler(context.Background(), newCallToolRequest(map[string]any{
		"modifier": 1,
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
}

// TestActionRollHandlerMapsRequestAndResponse ensures inputs and outputs are mapped consistently.
func TestActionRollHandlerMapsRequestAndResponse(t *testing.T) {
	difficulty := int32(7)
	client := &fakeDiceRollClient{
		response: &pb.ActionRollResponse{
			Duality: &pb.DualityDice{
				HopeD12: 4,
				FearD12: 6,
			},
			Total:      17,
			Difficulty: &difficulty,
			Outcome:    pb.Outcome_SUCCESS_WITH_HOPE,
		},
	}
	handler := actionRollHandler(client)

	result, err := handler(context.Background(), newCallToolRequest(map[string]any{
		"modifier":   7,
		"difficulty": 7,
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || result.IsError {
		t.Fatal("expected success result")
	}
	if client.lastRequest == nil {
		t.Fatal("expected gRPC request")
	}
	if client.lastRequest.Modifier != 7 {
		t.Fatalf("expected modifier 7, got %d", client.lastRequest.Modifier)
	}
	if client.lastRequest.Difficulty == nil || *client.lastRequest.Difficulty != 7 {
		t.Fatalf("expected difficulty 7, got %v", client.lastRequest.Difficulty)
	}

	structured, ok := result.StructuredContent.(ActionRollResult)
	if !ok {
		t.Fatalf("expected ActionRollResult, got %T", result.StructuredContent)
	}
	if structured.Hope != 4 || structured.Fear != 6 || structured.Total != 17 {
		t.Fatalf("unexpected dice output: %+v", structured)
	}
	if structured.Modifier != 7 {
		t.Fatalf("expected modifier 7, got %d", structured.Modifier)
	}
	if structured.Outcome != pb.Outcome_SUCCESS_WITH_HOPE.String() {
		t.Fatalf("expected outcome %q, got %q", pb.Outcome_SUCCESS_WITH_HOPE.String(), structured.Outcome)
	}
	if structured.Difficulty == nil || *structured.Difficulty != 7 {
		t.Fatalf("expected difficulty 7, got %v", structured.Difficulty)
	}
}

// TestRollDiceHandlerRejectsMissingDice ensures empty dice requests return an error result.
func TestRollDiceHandlerRejectsMissingDice(t *testing.T) {
	client := &fakeDiceRollClient{}
	handler := rollDiceHandler(client)

	result, err := handler(context.Background(), newRollDiceCallToolRequest(map[string]any{
		"dice": []map[string]any{},
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
	if client.lastRollDiceRequest != nil {
		t.Fatal("expected no gRPC call on invalid input")
	}
}

// TestRollDiceHandlerRejectsInvalidDice ensures invalid dice specs return an error result.
func TestRollDiceHandlerRejectsInvalidDice(t *testing.T) {
	client := &fakeDiceRollClient{}
	handler := rollDiceHandler(client)

	result, err := handler(context.Background(), newRollDiceCallToolRequest(map[string]any{
		"dice": []map[string]any{{"sides": -1, "count": 2}},
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
	if client.lastRollDiceRequest != nil {
		t.Fatal("expected no gRPC call on invalid input")
	}
}

// TestRollDiceHandlerReturnsClientError ensures gRPC errors are returned as tool errors.
func TestRollDiceHandlerReturnsClientError(t *testing.T) {
	client := &fakeDiceRollClient{rollDiceErr: errors.New("boom")}
	handler := rollDiceHandler(client)

	result, err := handler(context.Background(), newRollDiceCallToolRequest(map[string]any{
		"dice": []map[string]any{{"sides": 6, "count": 1}},
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error result")
	}
}

// TestRollDiceHandlerMapsRequestAndResponse ensures inputs and outputs are mapped consistently.
func TestRollDiceHandlerMapsRequestAndResponse(t *testing.T) {
	client := &fakeDiceRollClient{rollDiceResponse: &pb.RollDiceResponse{
		Rolls: []*pb.DiceRoll{
			{Sides: 6, Results: []int32{2, 5}, Total: 7},
			{Sides: 8, Results: []int32{4}, Total: 4},
		},
		Total: 11,
	}}

	handler := rollDiceHandler(client)

	result, err := handler(context.Background(), newRollDiceCallToolRequest(map[string]any{
		"dice": []map[string]any{{"sides": 6, "count": 2}, {"sides": 8, "count": 1}},
	}))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if result == nil || result.IsError {
		t.Fatal("expected success result")
	}
	if client.lastRollDiceRequest == nil {
		t.Fatal("expected gRPC request")
	}
	if len(client.lastRollDiceRequest.Dice) != 2 {
		t.Fatalf("expected 2 dice specs, got %d", len(client.lastRollDiceRequest.Dice))
	}
	if client.lastRollDiceRequest.Dice[0].Sides != 6 || client.lastRollDiceRequest.Dice[0].Count != 2 {
		t.Fatalf("unexpected first dice spec: %+v", client.lastRollDiceRequest.Dice[0])
	}
	if client.lastRollDiceRequest.Dice[1].Sides != 8 || client.lastRollDiceRequest.Dice[1].Count != 1 {
		t.Fatalf("unexpected second dice spec: %+v", client.lastRollDiceRequest.Dice[1])
	}

	structured, ok := result.StructuredContent.(RollDiceResult)
	if !ok {
		t.Fatalf("expected RollDiceResult, got %T", result.StructuredContent)
	}
	if structured.Total != 11 {
		t.Fatalf("expected total 11, got %d", structured.Total)
	}
	if len(structured.Rolls) != 2 {
		t.Fatalf("expected 2 rolls, got %d", len(structured.Rolls))
	}
	if structured.Rolls[0].Sides != 6 || structured.Rolls[0].Total != 7 {
		t.Fatalf("unexpected first roll: %+v", structured.Rolls[0])
	}
	if structured.Rolls[1].Sides != 8 || structured.Rolls[1].Total != 4 {
		t.Fatalf("unexpected second roll: %+v", structured.Rolls[1])
	}
}
