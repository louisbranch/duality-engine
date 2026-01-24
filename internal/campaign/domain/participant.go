package domain

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ParticipantRole describes the role of a participant in a campaign.
type ParticipantRole int

const (
	// ParticipantRoleUnspecified represents an invalid participant role value.
	ParticipantRoleUnspecified ParticipantRole = iota
	// ParticipantRoleGM indicates a game master.
	ParticipantRoleGM
	// ParticipantRolePlayer indicates a player.
	ParticipantRolePlayer
)

// Controller describes how a participant is controlled.
type Controller int

const (
	// ControllerUnspecified represents an invalid controller value.
	ControllerUnspecified Controller = iota
	// ControllerHuman indicates a human controller.
	ControllerHuman
	// ControllerAI indicates an AI controller.
	ControllerAI
)

var (
	// ErrEmptyDisplayName indicates a missing participant display name.
	ErrEmptyDisplayName = errors.New("display name is required")
	// ErrInvalidParticipantRole indicates a missing or invalid participant role.
	ErrInvalidParticipantRole = errors.New("participant role is required")
	// ErrEmptyCampaignID indicates a missing campaign ID.
	ErrEmptyCampaignID = errors.New("campaign id is required")
)

// Participant represents a participant (GM or player) in a campaign.
type Participant struct {
	ID          string
	CampaignID  string
	DisplayName string
	Role        ParticipantRole
	Controller  Controller
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateParticipantInput describes the metadata needed to create a participant.
type CreateParticipantInput struct {
	CampaignID  string
	DisplayName string
	Role        ParticipantRole
	Controller  Controller
}

// CreateParticipant creates a new participant with a generated ID and timestamps.
func CreateParticipant(input CreateParticipantInput, now func() time.Time, idGenerator func() (string, error)) (Participant, error) {
	if now == nil {
		now = time.Now
	}
	if idGenerator == nil {
		idGenerator = NewParticipantID
	}

	normalized, err := NormalizeCreateParticipantInput(input)
	if err != nil {
		return Participant{}, err
	}

	participantID, err := idGenerator()
	if err != nil {
		return Participant{}, fmt.Errorf("generate participant id: %w", err)
	}

	createdAt := now().UTC()
	return Participant{
		ID:          participantID,
		CampaignID:  normalized.CampaignID,
		DisplayName: normalized.DisplayName,
		Role:        normalized.Role,
		Controller:  normalized.Controller,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}, nil
}

// NormalizeCreateParticipantInput trims and validates participant input metadata.
func NormalizeCreateParticipantInput(input CreateParticipantInput) (CreateParticipantInput, error) {
	input.CampaignID = strings.TrimSpace(input.CampaignID)
	if input.CampaignID == "" {
		return CreateParticipantInput{}, ErrEmptyCampaignID
	}
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.DisplayName == "" {
		return CreateParticipantInput{}, ErrEmptyDisplayName
	}
	if input.Role == ParticipantRoleUnspecified {
		return CreateParticipantInput{}, ErrInvalidParticipantRole
	}
	if input.Controller == ControllerUnspecified {
		input.Controller = ControllerHuman
	}
	return input, nil
}

// NewParticipantID generates a URL-safe participant identifier.
func NewParticipantID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	// RFC 4122 variant and version bits for a v4 UUID.
	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw[:])
	return strings.ToLower(encoded), nil
}
