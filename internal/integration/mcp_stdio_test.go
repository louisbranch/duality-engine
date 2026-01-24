//go:build integration

// Package integration runs end-to-end integration tests.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/louisbranch/duality-engine/internal/app/server"
	dualitydomain "github.com/louisbranch/duality-engine/internal/duality/domain"
	"github.com/louisbranch/duality-engine/internal/mcp/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// TestMCPStdioEndToEnd validates MCP stdio integration end-to-end.
func TestMCPStdioEndToEnd(t *testing.T) {
	grpcAddr, stopServer := startGRPCServer(t)
	defer stopServer()

	clientSession, closeClient := startMCPClient(t, grpcAddr)
	defer closeClient()

	t.Run("list tools", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout())
		defer cancel()

		result, err := clientSession.ListTools(ctx, &mcp.ListToolsParams{})
		if err != nil {
			t.Fatalf("list tools: %v", err)
		}
		if result == nil {
			t.Fatal("list tools returned nil result")
		}

		toolNames := make(map[string]bool)
		for _, tool := range result.Tools {
			toolNames[tool.Name] = true
		}

		expected := []string{
			"duality_action_roll",
			"duality_outcome",
			"duality_rules_version",
			"roll_dice",
			"campaign_create",
		}
		for _, name := range expected {
			if !toolNames[name] {
				t.Fatalf("expected tool %q", name)
			}
		}
	})

	t.Run("duality outcome", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout())
		defer cancel()

		params := &mcp.CallToolParams{
			Name: "duality_outcome",
			Arguments: map[string]any{
				"hope":       10,
				"fear":       4,
				"modifier":   1,
				"difficulty": 10,
			},
		}
		result, err := clientSession.CallTool(ctx, params)
		if err != nil {
			t.Fatalf("call duality_outcome: %v", err)
		}
		if result == nil {
			t.Fatal("call duality_outcome returned nil")
		}
		if result.IsError {
			t.Fatalf("call duality_outcome returned error content: %+v", result.Content)
		}

		output := decodeStructuredContent[domain.DualityOutcomeResult](t, result.StructuredContent)
		if output.Hope != 10 || output.Fear != 4 || output.Modifier != 1 {
			t.Fatalf("unexpected dice output: %+v", output)
		}
		if output.Total != 15 {
			t.Fatalf("expected total 15, got %d", output.Total)
		}
		if output.Difficulty == nil || *output.Difficulty != 10 {
			t.Fatalf("expected difficulty 10, got %v", output.Difficulty)
		}
		if output.Outcome != "SUCCESS_WITH_HOPE" {
			t.Fatalf("expected outcome SUCCESS_WITH_HOPE, got %q", output.Outcome)
		}

		expected, err := dualitydomain.EvaluateOutcome(dualitydomain.OutcomeRequest{
			Hope:       10,
			Fear:       4,
			Modifier:   1,
			Difficulty: intPointer(10),
		})
		if err != nil {
			t.Fatalf("evaluate outcome: %v", err)
		}
		if output.Total != expected.Total || output.MeetsDifficulty != expected.MeetsDifficulty {
			t.Fatalf("unexpected outcome totals: %+v", output)
		}
	})

	t.Run("campaign create and list", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout())
		defer cancel()

		params := &mcp.CallToolParams{
			Name: "campaign_create",
			Arguments: map[string]any{
				"name":         "Stormbound",
				"gm_mode":      "HUMAN",
				"player_slots": 4,
				"theme_prompt": "sea and thunder",
			},
		}
		result, err := clientSession.CallTool(ctx, params)
		if err != nil {
			t.Fatalf("call campaign_create: %v", err)
		}
		if result == nil {
			t.Fatal("call campaign_create returned nil")
		}
		if result.IsError {
			t.Fatalf("campaign_create returned error content: %+v", result.Content)
		}

		output := decodeStructuredContent[domain.CampaignCreateResult](t, result.StructuredContent)
		if output.ID == "" {
			t.Fatal("campaign_create returned empty id")
		}
		if output.Name != "Stormbound" {
			t.Fatalf("expected campaign name Stormbound, got %q", output.Name)
		}
		if output.GmMode != "HUMAN" {
			t.Fatalf("expected gm_mode HUMAN, got %q", output.GmMode)
		}

		resource, err := clientSession.ReadResource(ctx, &mcp.ReadResourceParams{URI: "campaigns://list"})
		if err != nil {
			t.Fatalf("read campaigns://list: %v", err)
		}
		if resource == nil || len(resource.Contents) == 0 {
			t.Fatalf("read campaigns://list returned no contents: %+v", resource)
		}

		payload := parseCampaignListPayload(t, resource.Contents[0].Text)
		entry, found := findCampaignByID(payload, output.ID)
		if !found {
			t.Fatalf("campaign %q not found in list", output.ID)
		}
		if entry.Name != output.Name {
			t.Fatalf("expected campaign name %q, got %q", output.Name, entry.Name)
		}
		if entry.GmMode != output.GmMode {
			t.Fatalf("expected gm_mode %q, got %q", output.GmMode, entry.GmMode)
		}
		createdAt := parseRFC3339(t, entry.CreatedAt)
		updatedAt := parseRFC3339(t, entry.UpdatedAt)
		if updatedAt.Before(createdAt) {
			t.Fatalf("expected updated_at after created_at: %v < %v", updatedAt, createdAt)
		}
	})

	t.Run("rules metadata", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), integrationTimeout())
		defer cancel()

		params := &mcp.CallToolParams{Name: "duality_rules_version"}
		result, err := clientSession.CallTool(ctx, params)
		if err != nil {
			t.Fatalf("call duality_rules_version: %v", err)
		}
		if result == nil {
			t.Fatal("call duality_rules_version returned nil")
		}
		if result.IsError {
			t.Fatalf("duality_rules_version returned error content: %+v", result.Content)
		}

		output := decodeStructuredContent[domain.RulesVersionResult](t, result.StructuredContent)
		expected := dualitydomain.RulesVersion()
		if output.System != expected.System {
			t.Fatalf("expected system %q, got %q", expected.System, output.System)
		}
		if output.Module != expected.Module {
			t.Fatalf("expected module %q, got %q", expected.Module, output.Module)
		}
		if output.RulesVersion != expected.RulesVersion {
			t.Fatalf("expected rules version %q, got %q", expected.RulesVersion, output.RulesVersion)
		}
		if output.DiceModel != expected.DiceModel {
			t.Fatalf("expected dice model %q, got %q", expected.DiceModel, output.DiceModel)
		}

		expectedOutcomes := []string{
			"ROLL_WITH_HOPE",
			"ROLL_WITH_FEAR",
			"SUCCESS_WITH_HOPE",
			"SUCCESS_WITH_FEAR",
			"FAILURE_WITH_HOPE",
			"FAILURE_WITH_FEAR",
			"CRITICAL_SUCCESS",
		}
		if len(output.Outcomes) != len(expectedOutcomes) {
			t.Fatalf("expected %d outcomes, got %d", len(expectedOutcomes), len(output.Outcomes))
		}
		for i, expectedOutcome := range expectedOutcomes {
			if output.Outcomes[i] != expectedOutcome {
				t.Fatalf("expected outcome %q at index %d, got %q", expectedOutcome, i, output.Outcomes[i])
			}
		}
	})
}

// integrationTimeout returns the default timeout for integration calls.
func integrationTimeout() time.Duration {
	return 10 * time.Second
}

// startGRPCServer boots the gRPC server and returns its address and shutdown function.
func startGRPCServer(t *testing.T) (string, func()) {
	t.Helper()

	setTempDBPath(t)

	ctx, cancel := context.WithCancel(context.Background())
	grpcServer, err := server.New(0)
	if err != nil {
		cancel()
		t.Fatalf("new gRPC server: %v", err)
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(ctx)
	}()

	addr := normalizeAddress(t, grpcServer.Addr())
	waitForGRPCHealth(t, addr)
	stop := func() {
		cancel()
		select {
		case err := <-serveErr:
			if err != nil {
				t.Fatalf("gRPC server error: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for gRPC server to stop")
		}
	}

	return addr, stop
}

// startMCPClient boots the MCP stdio process and returns a client session and shutdown function.
func startMCPClient(t *testing.T, grpcAddr string) (*mcp.ClientSession, func()) {
	t.Helper()

	cmd := exec.Command("go", "run", "./cmd/mcp")
	cmd.Dir = repoRoot(t)
	cmd.Env = append(os.Environ(), fmt.Sprintf("DUALITY_GRPC_ADDR=%s", grpcAddr))
	cmd.Stderr = os.Stderr

	transport := &mcp.CommandTransport{Command: cmd}
	client := mcp.NewClient(&mcp.Implementation{Name: "integration-client", Version: "dev"}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientSession, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}

	closeClient := func() {
		closeErr := clientSession.Close()
		if closeErr != nil {
			t.Fatalf("close MCP client: %v", closeErr)
		}
	}

	return clientSession, closeClient
}

// decodeStructuredContent decodes structured MCP content into the target type.
func decodeStructuredContent[T any](t *testing.T, value any) T {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var output T
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("unmarshal structured content: %v", err)
	}
	return output
}

// parseCampaignListPayload decodes a campaign list JSON payload.
func parseCampaignListPayload(t *testing.T, raw string) domain.CampaignListPayload {
	t.Helper()

	var payload domain.CampaignListPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("unmarshal campaign list payload: %v", err)
	}
	return payload
}

// findCampaignByID searches for a campaign entry by ID.
func findCampaignByID(payload domain.CampaignListPayload, id string) (domain.CampaignListEntry, bool) {
	for _, campaign := range payload.Campaigns {
		if campaign.ID == id {
			return campaign, true
		}
	}
	return domain.CampaignListEntry{}, false
}

// parseRFC3339 parses an RFC3339 timestamp string.
func parseRFC3339(t *testing.T, value string) time.Time {
	t.Helper()

	if value == "" {
		t.Fatal("expected non-empty timestamp")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("parse timestamp %q: %v", value, err)
	}
	return parsed
}

// setTempDBPath configures a temporary database for integration tests.
func setTempDBPath(t *testing.T) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "duality.db")
	t.Setenv("DUALITY_DB_PATH", path)
}

// repoRoot returns the repository root by walking up to go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve runtime caller")
	}

	dir := filepath.Dir(filename)
	for {
		candidate := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatalf("go.mod not found from %s", filename)
	return ""
}

// normalizeAddress maps wildcard listener hosts to localhost.
func normalizeAddress(t *testing.T, addr string) string {
	t.Helper()

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split address %q: %v", addr, err)
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port)
}

// waitForGRPCHealth waits for the gRPC health check to report SERVING.
func waitForGRPCHealth(t *testing.T, addr string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		t.Fatalf("dial gRPC server: %v", err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	backoff := 100 * time.Millisecond
	for {
		callCtx, callCancel := context.WithTimeout(ctx, time.Second)
		response, err := healthClient.Check(callCtx, &grpc_health_v1.HealthCheckRequest{Service: ""})
		callCancel()
		if err == nil && response.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
			return
		}

		select {
		case <-ctx.Done():
			if err != nil {
				t.Fatalf("wait for gRPC health: %v", err)
			}
			t.Fatalf("wait for gRPC health: %v", ctx.Err())
		case <-time.After(backoff):
		}

		if backoff < time.Second {
			backoff *= 2
			if backoff > time.Second {
				backoff = time.Second
			}
		}
	}
}

// intPointer returns a pointer to the provided int value.
func intPointer(value int) *int {
	return &value
}
