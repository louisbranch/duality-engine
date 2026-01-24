package domain

import (
	"encoding/base32"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCreateParticipantNormalizesInput(t *testing.T) {
	fixedTime := time.Date(2026, 1, 23, 10, 0, 0, 0, time.UTC)
	input := CreateParticipantInput{
		CampaignID:  "camp-123",
		DisplayName: "  Alice  ",
		Role:        ParticipantRolePlayer,
		Controller:  ControllerHuman,
	}

	participant, err := CreateParticipant(input, func() time.Time { return fixedTime }, func() (string, error) {
		return "part-456", nil
	})
	if err != nil {
		t.Fatalf("create participant: %v", err)
	}

	if participant.ID != "part-456" {
		t.Fatalf("expected id part-456, got %q", participant.ID)
	}
	if participant.CampaignID != "camp-123" {
		t.Fatalf("expected campaign id camp-123, got %q", participant.CampaignID)
	}
	if participant.DisplayName != "Alice" {
		t.Fatalf("expected trimmed display name, got %q", participant.DisplayName)
	}
	if participant.Role != ParticipantRolePlayer {
		t.Fatalf("expected role player, got %v", participant.Role)
	}
	if participant.Controller != ControllerHuman {
		t.Fatalf("expected controller human, got %v", participant.Controller)
	}
	if !participant.CreatedAt.Equal(fixedTime) || !participant.UpdatedAt.Equal(fixedTime) {
		t.Fatalf("expected timestamps to match fixed time")
	}
}

func TestCreateParticipantDefaultsController(t *testing.T) {
	fixedTime := time.Date(2026, 1, 23, 10, 0, 0, 0, time.UTC)
	input := CreateParticipantInput{
		CampaignID:  "camp-123",
		DisplayName: "Bob",
		Role:        ParticipantRoleGM,
		Controller:  ControllerUnspecified,
	}

	participant, err := CreateParticipant(input, func() time.Time { return fixedTime }, func() (string, error) {
		return "part-789", nil
	})
	if err != nil {
		t.Fatalf("create participant: %v", err)
	}

	if participant.Controller != ControllerHuman {
		t.Fatalf("expected default controller human, got %v", participant.Controller)
	}
}

func TestNormalizeCreateParticipantInputValidation(t *testing.T) {
	tests := []struct {
		name  string
		input CreateParticipantInput
		err   error
	}{
		{
			name: "empty campaign id",
			input: CreateParticipantInput{
				CampaignID:  "   ",
				DisplayName: "Alice",
				Role:        ParticipantRolePlayer,
			},
			err: ErrEmptyCampaignID,
		},
		{
			name: "empty display name",
			input: CreateParticipantInput{
				CampaignID:  "camp-123",
				DisplayName: "   ",
				Role:        ParticipantRolePlayer,
			},
			err: ErrEmptyDisplayName,
		},
		{
			name: "missing role",
			input: CreateParticipantInput{
				CampaignID:  "camp-123",
				DisplayName: "Alice",
				Role:        ParticipantRoleUnspecified,
			},
			err: ErrInvalidParticipantRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeCreateParticipantInput(tt.input)
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error %v, got %v", tt.err, err)
			}
		})
	}
}

func TestNewParticipantIDFormat(t *testing.T) {
	id, err := NewParticipantID()
	if err != nil {
		t.Fatalf("new participant id: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
	if strings.Contains(id, "=") {
		t.Fatal("expected no padding")
	}
	if len(id) != 26 {
		t.Fatalf("expected 26-character id, got %d", len(id))
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '2' || r > '7') {
			t.Fatalf("unexpected character %q in id", r)
		}
	}

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(id))
	if err != nil {
		t.Fatalf("decode id: %v", err)
	}
	if len(decoded) != 16 {
		t.Fatalf("expected 16 decoded bytes, got %d", len(decoded))
	}
}
