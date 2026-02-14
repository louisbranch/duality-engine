package domain

import (
	"context"
	"encoding/json"
	"testing"

	statev1 "github.com/louisbranch/fracturing.space/api/gen/go/game/v1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSetContextTool(t *testing.T) {
	tool := SetContextTool()
	if tool == nil {
		t.Fatal("expected non-nil tool")
	}
	if tool.Name != "set_context" {
		t.Errorf("expected name %q, got %q", "set_context", tool.Name)
	}
}

func TestContextResource(t *testing.T) {
	res := ContextResource()
	if res == nil {
		t.Fatal("expected non-nil resource")
	}
	if res.URI != "context://current" {
		t.Errorf("expected URI %q, got %q", "context://current", res.URI)
	}
	if res.MIMEType != "application/json" {
		t.Errorf("expected MIME type %q, got %q", "application/json", res.MIMEType)
	}
}

func TestContextResourceHandler(t *testing.T) {
	t.Run("nil getter returns error", func(t *testing.T) {
		handler := ContextResourceHandler(nil)
		_, err := handler(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: "context://current"},
		})
		if err == nil {
			t.Fatal("expected error for nil getter")
		}
	})

	t.Run("wrong URI returns error", func(t *testing.T) {
		handler := ContextResourceHandler(func() Context { return Context{} })
		_, err := handler(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: "context://wrong"},
		})
		if err == nil {
			t.Fatal("expected error for wrong URI")
		}
	})

	t.Run("empty context", func(t *testing.T) {
		handler := ContextResourceHandler(func() Context { return Context{} })
		result, err := handler(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: "context://current"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Contents) != 1 {
			t.Fatalf("expected 1 content, got %d", len(result.Contents))
		}

		var payload ContextResourcePayload
		if err := json.Unmarshal([]byte(result.Contents[0].Text), &payload); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if payload.Context.CampaignID != nil {
			t.Error("expected nil campaign_id")
		}
	})

	t.Run("populated context", func(t *testing.T) {
		handler := ContextResourceHandler(func() Context {
			return Context{CampaignID: "camp-1", SessionID: "sess-1", ParticipantID: "part-1"}
		})
		result, err := handler(context.Background(), &mcp.ReadResourceRequest{
			Params: &mcp.ReadResourceParams{URI: "context://current"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var payload ContextResourcePayload
		if err := json.Unmarshal([]byte(result.Contents[0].Text), &payload); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if payload.Context.CampaignID == nil || *payload.Context.CampaignID != "camp-1" {
			t.Errorf("expected campaign_id %q, got %v", "camp-1", payload.Context.CampaignID)
		}
		if payload.Context.SessionID == nil || *payload.Context.SessionID != "sess-1" {
			t.Errorf("expected session_id %q, got %v", "sess-1", payload.Context.SessionID)
		}
	})

	t.Run("nil request uses default URI", func(t *testing.T) {
		handler := ContextResourceHandler(func() Context { return Context{} })
		result, err := handler(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Contents[0].URI != "context://current" {
			t.Errorf("expected default URI, got %q", result.Contents[0].URI)
		}
	})
}

func TestSetContextHandler(t *testing.T) {
	t.Run("empty campaign_id returns error", func(t *testing.T) {
		handler := SetContextHandler(nil, nil, nil, func(_ Context) {}, func() Context { return Context{} }, nil)
		_, _, err := handler(context.Background(), nil, SetContextInput{CampaignID: ""})
		if err == nil {
			t.Fatal("expected error for empty campaign_id")
		}
	})

	t.Run("whitespace campaign_id returns error", func(t *testing.T) {
		handler := SetContextHandler(nil, nil, nil, func(_ Context) {}, func() Context { return Context{} }, nil)
		_, _, err := handler(context.Background(), nil, SetContextInput{CampaignID: "  "})
		if err == nil {
			t.Fatal("expected error for whitespace campaign_id")
		}
	})

	t.Run("campaign not found", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getErr: status.Error(codes.NotFound, "not found"),
		}
		handler := SetContextHandler(campaignClient, nil, nil, func(_ Context) {}, func() Context { return Context{} }, nil)
		_, _, err := handler(context.Background(), nil, SetContextInput{CampaignID: "camp-1"})
		if err == nil {
			t.Fatal("expected error for not found campaign")
		}
	})

	t.Run("success with campaign only", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		var setCtx Context
		handler := SetContextHandler(
			campaignClient, nil, nil,
			func(c Context) { setCtx = c },
			func() Context { return setCtx },
			nil,
		)
		toolResult, result, err := handler(context.Background(), nil, SetContextInput{CampaignID: "camp-1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Context.CampaignID != "camp-1" {
			t.Errorf("expected campaign_id %q, got %q", "camp-1", result.Context.CampaignID)
		}
		if toolResult == nil {
			t.Fatal("expected non-nil tool result")
		}
		if setCtx.CampaignID != "camp-1" {
			t.Errorf("expected setContext called with camp-1, got %q", setCtx.CampaignID)
		}
	})

	t.Run("success with session and participant", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		sessionClient := &fakeSessionClient{
			getResp: &statev1.GetSessionResponse{
				Session: testSession("sess-1", "camp-1", "Session 1", statev1.SessionStatus_SESSION_ACTIVE),
			},
		}
		participantClient := &fakeParticipantClient{
			getResp: &statev1.GetParticipantResponse{
				Participant: testParticipant("part-1", "camp-1", "Alice", statev1.ParticipantRole_PLAYER),
			},
		}

		var setCtx Context
		handler := SetContextHandler(
			campaignClient, sessionClient, participantClient,
			func(c Context) { setCtx = c },
			func() Context { return setCtx },
			nil,
		)
		_, result, err := handler(context.Background(), nil, SetContextInput{
			CampaignID:    "camp-1",
			SessionID:     "sess-1",
			ParticipantID: "part-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Context.SessionID != "sess-1" {
			t.Errorf("expected session_id %q, got %q", "sess-1", result.Context.SessionID)
		}
		if result.Context.ParticipantID != "part-1" {
			t.Errorf("expected participant_id %q, got %q", "part-1", result.Context.ParticipantID)
		}
	})

	t.Run("whitespace session_id treated as omitted", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		var setCtx Context
		handler := SetContextHandler(
			campaignClient, nil, nil,
			func(c Context) { setCtx = c },
			func() Context { return setCtx },
			nil,
		)
		_, result, err := handler(context.Background(), nil, SetContextInput{
			CampaignID: "camp-1",
			SessionID:  "  ",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Context.SessionID != "" {
			t.Errorf("expected empty session_id, got %q", result.Context.SessionID)
		}
	})

	t.Run("session not found", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		sessionClient := &fakeSessionClient{
			getErr: status.Error(codes.NotFound, "not found"),
		}
		handler := SetContextHandler(
			campaignClient, sessionClient, nil,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID: "camp-1",
			SessionID:  "sess-bad",
		})
		if err == nil {
			t.Fatal("expected error for not found session")
		}
	})

	t.Run("participant not found", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		participantClient := &fakeParticipantClient{
			getErr: status.Error(codes.NotFound, "not found"),
		}
		handler := SetContextHandler(
			campaignClient, nil, participantClient,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID:    "camp-1",
			ParticipantID: "part-bad",
		})
		if err == nil {
			t.Fatal("expected error for not found participant")
		}
	})

	t.Run("session invalid argument", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		sessionClient := &fakeSessionClient{
			getErr: status.Error(codes.InvalidArgument, "session does not belong to campaign"),
		}
		handler := SetContextHandler(
			campaignClient, sessionClient, nil,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID: "camp-1",
			SessionID:  "sess-wrong-camp",
		})
		if err == nil {
			t.Fatal("expected error for invalid argument session")
		}
	})

	t.Run("participant invalid argument", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		participantClient := &fakeParticipantClient{
			getErr: status.Error(codes.InvalidArgument, "participant does not belong to campaign"),
		}
		handler := SetContextHandler(
			campaignClient, nil, participantClient,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID:    "camp-1",
			ParticipantID: "part-wrong-camp",
		})
		if err == nil {
			t.Fatal("expected error for invalid argument participant")
		}
	})

	t.Run("campaign gRPC error", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getErr: status.Error(codes.Internal, "internal error"),
		}
		handler := SetContextHandler(
			campaignClient, nil, nil,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID: "camp-1",
		})
		if err == nil {
			t.Fatal("expected error for gRPC internal error")
		}
	})

	t.Run("session gRPC error", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		sessionClient := &fakeSessionClient{
			getErr: status.Error(codes.Internal, "internal error"),
		}
		handler := SetContextHandler(
			campaignClient, sessionClient, nil,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID: "camp-1",
			SessionID:  "sess-1",
		})
		if err == nil {
			t.Fatal("expected error for session gRPC internal error")
		}
	})

	t.Run("participant gRPC error", func(t *testing.T) {
		campaignClient := &fakeCampaignClient{
			getResp: &statev1.GetCampaignResponse{
				Campaign: testCampaign("camp-1", "Test", statev1.CampaignStatus_ACTIVE),
			},
		}
		participantClient := &fakeParticipantClient{
			getErr: status.Error(codes.Internal, "internal error"),
		}
		handler := SetContextHandler(
			campaignClient, nil, participantClient,
			func(_ Context) {},
			func() Context { return Context{} },
			nil,
		)
		_, _, err := handler(context.Background(), nil, SetContextInput{
			CampaignID:    "camp-1",
			ParticipantID: "part-1",
		})
		if err == nil {
			t.Fatal("expected error for participant gRPC internal error")
		}
	})
}
