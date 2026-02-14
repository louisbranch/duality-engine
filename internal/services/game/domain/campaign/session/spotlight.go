package session

import (
	"fmt"
	"strings"
	"time"
)

// SpotlightType describes who currently has the spotlight.
type SpotlightType string

const (
	SpotlightTypeGM        SpotlightType = "gm"
	SpotlightTypeCharacter SpotlightType = "character"
)

// Spotlight tracks the current session focus.
type Spotlight struct {
	CampaignID         string
	SessionID          string
	Type               SpotlightType
	CharacterID        string
	UpdatedAt          time.Time
	UpdatedByActorType string
	UpdatedByActorID   string
}

// NormalizeSpotlightType validates and normalizes a spotlight type value.
func NormalizeSpotlightType(value string) (SpotlightType, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	if trimmed == "" {
		return "", fmt.Errorf("spotlight type is required")
	}
	switch SpotlightType(trimmed) {
	case SpotlightTypeGM, SpotlightTypeCharacter:
		return SpotlightType(trimmed), nil
	default:
		return "", fmt.Errorf("spotlight type %q is invalid", value)
	}
}

// ValidateSpotlightTarget enforces target requirements based on spotlight type.
func ValidateSpotlightTarget(spotlightType SpotlightType, characterID string) error {
	characterID = strings.TrimSpace(characterID)
	if spotlightType == SpotlightTypeCharacter && characterID == "" {
		return fmt.Errorf("character id is required for character spotlight")
	}
	if spotlightType == SpotlightTypeGM && characterID != "" {
		return fmt.Errorf("character id must be empty for gm spotlight")
	}
	return nil
}
