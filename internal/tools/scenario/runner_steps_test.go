package scenario

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	gamev1 "github.com/louisbranch/fracturing.space/api/gen/go/game/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestRunParticipantStepDefaults(t *testing.T) {
	var gotRequest *gamev1.CreateParticipantRequest
	participantClient := &fakeParticipantClient{
		create: func(_ context.Context, req *gamev1.CreateParticipantRequest, _ ...grpc.CallOption) (*gamev1.CreateParticipantResponse, error) {
			gotRequest = req
			return &gamev1.CreateParticipantResponse{
				Participant: &gamev1.Participant{Id: "participant-1"},
			}, nil
		},
	}

	runner := &Runner{
		assertions: Assertions{Mode: AssertionStrict},
		env: scenarioEnv{
			participantClient: participantClient,
			eventClient:       &fakeEventClient{},
		},
	}
	state := &scenarioState{campaignID: "campaign-1", participants: map[string]string{}}
	step := Step{Kind: "participant", Args: map[string]any{"name": "Alice"}}

	if err := runner.runParticipantStep(context.Background(), state, step); err != nil {
		t.Fatalf("runParticipantStep: %v", err)
	}
	if gotRequest == nil {
		t.Fatal("expected create participant request")
	}
	if gotRequest.GetRole() != gamev1.ParticipantRole_PLAYER {
		t.Fatalf("role = %s, want PLAYER", gotRequest.GetRole().String())
	}
	if gotRequest.GetController() != gamev1.Controller_CONTROLLER_HUMAN {
		t.Fatalf("controller = %s, want HUMAN", gotRequest.GetController().String())
	}
}

func TestRunCharacterStepControlParticipant(t *testing.T) {
	var controlRequest *gamev1.SetDefaultControlRequest
	characterClient := &fakeCharacterClient{
		create: func(_ context.Context, req *gamev1.CreateCharacterRequest, _ ...grpc.CallOption) (*gamev1.CreateCharacterResponse, error) {
			return &gamev1.CreateCharacterResponse{
				Character: &gamev1.Character{Id: "character-1"},
			}, nil
		},
		patchProfile: func(context.Context, *gamev1.PatchCharacterProfileRequest, ...grpc.CallOption) (*gamev1.PatchCharacterProfileResponse, error) {
			return &gamev1.PatchCharacterProfileResponse{}, nil
		},
		setDefaultControl: func(_ context.Context, req *gamev1.SetDefaultControlRequest, _ ...grpc.CallOption) (*gamev1.SetDefaultControlResponse, error) {
			controlRequest = req
			return &gamev1.SetDefaultControlResponse{}, nil
		},
	}

	runner := &Runner{
		assertions: Assertions{Mode: AssertionStrict},
		env: scenarioEnv{
			characterClient: characterClient,
			snapshotClient:  &fakeSnapshotClient{},
			eventClient:     &fakeEventClient{},
		},
	}
	state := &scenarioState{
		campaignID:         "campaign-1",
		ownerParticipantID: "owner-1",
		participants:       map[string]string{"John": "participant-1"},
		actors:             map[string]string{},
	}
	step := Step{Kind: "character", Args: map[string]any{
		"name":        "Frodo",
		"control":     "participant",
		"participant": "John",
	}}

	if err := runner.runCharacterStep(context.Background(), state, step); err != nil {
		t.Fatalf("runCharacterStep: %v", err)
	}
	if controlRequest == nil {
		t.Fatal("expected SetDefaultControl request")
	}
	if got := controlRequest.GetParticipantId(); got == nil || got.GetValue() != "participant-1" {
		t.Fatalf("participant_id = %v, want participant-1", got)
	}
}

func TestRunCharacterStepControlGM(t *testing.T) {
	var controlRequest *gamev1.SetDefaultControlRequest
	characterClient := &fakeCharacterClient{
		create: func(_ context.Context, req *gamev1.CreateCharacterRequest, _ ...grpc.CallOption) (*gamev1.CreateCharacterResponse, error) {
			return &gamev1.CreateCharacterResponse{
				Character: &gamev1.Character{Id: "character-1"},
			}, nil
		},
		patchProfile: func(context.Context, *gamev1.PatchCharacterProfileRequest, ...grpc.CallOption) (*gamev1.PatchCharacterProfileResponse, error) {
			return &gamev1.PatchCharacterProfileResponse{}, nil
		},
		setDefaultControl: func(_ context.Context, req *gamev1.SetDefaultControlRequest, _ ...grpc.CallOption) (*gamev1.SetDefaultControlResponse, error) {
			controlRequest = req
			return &gamev1.SetDefaultControlResponse{}, nil
		},
	}

	runner := &Runner{
		assertions: Assertions{Mode: AssertionStrict},
		env: scenarioEnv{
			characterClient: characterClient,
			snapshotClient:  &fakeSnapshotClient{},
			eventClient:     &fakeEventClient{},
		},
	}
	state := &scenarioState{
		campaignID:         "campaign-1",
		ownerParticipantID: "owner-1",
		actors:             map[string]string{},
	}
	step := Step{Kind: "character", Args: map[string]any{
		"name":    "Frodo",
		"control": "gm",
	}}

	if err := runner.runCharacterStep(context.Background(), state, step); err != nil {
		t.Fatalf("runCharacterStep: %v", err)
	}
	if controlRequest == nil {
		t.Fatal("expected SetDefaultControl request")
	}
	if controlRequest.GetParticipantId() != nil {
		t.Fatalf("participant_id = %v, want nil", controlRequest.GetParticipantId())
	}
}

func TestRunCampaignStepVerboseLogging(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := log.New(buffer, "", 0)

	runner := &Runner{
		assertions: Assertions{Mode: AssertionStrict},
		logger:     logger,
		verbose:    true,
		env: scenarioEnv{
			campaignClient: &fakeCampaignClient{
				create: func(_ context.Context, req *gamev1.CreateCampaignRequest, _ ...grpc.CallOption) (*gamev1.CreateCampaignResponse, error) {
					return &gamev1.CreateCampaignResponse{
						Campaign:         &gamev1.Campaign{Id: "campaign-1"},
						OwnerParticipant: &gamev1.Participant{Id: "participant-1"},
					}, nil
				},
			},
			eventClient: &fakeEventClient{},
		},
	}

	state := &scenarioState{}
	step := Step{Kind: "campaign", Args: map[string]any{"name": "Test", "system": "DAGGERHEART"}}
	if err := runner.runCampaignStep(context.Background(), state, step); err != nil {
		t.Fatalf("runCampaignStep: %v", err)
	}
	if !strings.Contains(buffer.String(), "campaign created") {
		t.Fatalf("expected verbose log to include campaign created")
	}
}

func TestRunCampaignStepRequiresOwnerParticipant(t *testing.T) {
	runner := &Runner{
		assertions: Assertions{Mode: AssertionStrict},
		env: scenarioEnv{
			campaignClient: &fakeCampaignClient{
				create: func(_ context.Context, req *gamev1.CreateCampaignRequest, _ ...grpc.CallOption) (*gamev1.CreateCampaignResponse, error) {
					return &gamev1.CreateCampaignResponse{
						Campaign: &gamev1.Campaign{Id: "campaign-1"},
					}, nil
				},
			},
			eventClient: &fakeEventClient{},
		},
	}

	state := &scenarioState{}
	step := Step{Kind: "campaign", Args: map[string]any{"name": "Test", "system": "DAGGERHEART"}}
	if err := runner.runCampaignStep(context.Background(), state, step); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseControlInvalid(t *testing.T) {
	_, err := parseControl("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseParticipantRoleController(t *testing.T) {
	role, err := parseParticipantRole("GM")
	if err != nil {
		t.Fatalf("parseParticipantRole: %v", err)
	}
	if role != gamev1.ParticipantRole_GM {
		t.Fatalf("role = %s, want GM", role.String())
	}
	controller, err := parseController("AI")
	if err != nil {
		t.Fatalf("parseController: %v", err)
	}
	if controller != gamev1.Controller_CONTROLLER_AI {
		t.Fatalf("controller = %s, want AI", controller.String())
	}
}

func TestSetDefaultControlRequestWithoutParticipant(t *testing.T) {
	request := &gamev1.SetDefaultControlRequest{}
	if request.GetParticipantId() != nil {
		t.Fatalf("participant_id = %v, want nil", request.GetParticipantId())
	}
	request.ParticipantId = wrapperspb.String("participant-1")
	if got := request.GetParticipantId().GetValue(); got != "participant-1" {
		t.Fatalf("participant_id = %s, want participant-1", got)
	}
}
