package session

import (
	"fmt"
	"strings"
)

// GateStatus describes the lifecycle state of a session gate.
type GateStatus string

const (
	GateStatusOpen      GateStatus = "open"
	GateStatusResolved  GateStatus = "resolved"
	GateStatusAbandoned GateStatus = "abandoned"
)

// NormalizeGateType validates and normalizes a gate type value.
func NormalizeGateType(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("gate type is required")
	}
	return strings.ToLower(trimmed), nil
}

// NormalizeGateReason trims a gate reason string.
func NormalizeGateReason(value string) string {
	return strings.TrimSpace(value)
}
