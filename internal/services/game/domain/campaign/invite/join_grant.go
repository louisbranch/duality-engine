package invite

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	apperrors "github.com/louisbranch/fracturing.space/internal/platform/errors"
)

const (
	EnvJoinGrantIssuer    = "FRACTURING_SPACE_JOIN_GRANT_ISSUER"
	EnvJoinGrantAudience  = "FRACTURING_SPACE_JOIN_GRANT_AUDIENCE"
	EnvJoinGrantPublicKey = "FRACTURING_SPACE_JOIN_GRANT_PUBLIC_KEY"
)

// JoinGrantConfig defines how join grants are verified.
type JoinGrantConfig struct {
	Issuer   string
	Audience string
	Key      ed25519.PublicKey
	Now      func() time.Time
}

// JoinGrantExpectation defines the expected identity for a join grant.
type JoinGrantExpectation struct {
	CampaignID string
	InviteID   string
	UserID     string
}

// JoinGrantClaims captures validated join grant claims.
type JoinGrantClaims struct {
	Issuer     string
	Audience   []string
	ExpiresAt  time.Time
	NotBefore  time.Time
	IssuedAt   time.Time
	JWTID      string
	CampaignID string
	InviteID   string
	UserID     string
}

// LoadJoinGrantConfigFromEnv reads join grant verification configuration.
func LoadJoinGrantConfigFromEnv(now func() time.Time) (JoinGrantConfig, error) {
	issuer := strings.TrimSpace(os.Getenv(EnvJoinGrantIssuer))
	if issuer == "" {
		return JoinGrantConfig{}, fmt.Errorf("%s is required", EnvJoinGrantIssuer)
	}
	audience := strings.TrimSpace(os.Getenv(EnvJoinGrantAudience))
	if audience == "" {
		return JoinGrantConfig{}, fmt.Errorf("%s is required", EnvJoinGrantAudience)
	}
	keyRaw := strings.TrimSpace(os.Getenv(EnvJoinGrantPublicKey))
	if keyRaw == "" {
		return JoinGrantConfig{}, fmt.Errorf("%s is required", EnvJoinGrantPublicKey)
	}
	keyBytes, err := decodeBase64(keyRaw)
	if err != nil {
		return JoinGrantConfig{}, fmt.Errorf("decode join grant public key: %w", err)
	}
	if len(keyBytes) != ed25519.PublicKeySize {
		return JoinGrantConfig{}, fmt.Errorf("join grant public key must be %d bytes", ed25519.PublicKeySize)
	}
	if now == nil {
		now = time.Now
	}
	return JoinGrantConfig{
		Issuer:   issuer,
		Audience: audience,
		Key:      ed25519.PublicKey(keyBytes),
		Now:      now,
	}, nil
}

// ValidateJoinGrant verifies a join grant token and validates expected claims.
func ValidateJoinGrant(grant string, expected JoinGrantExpectation, cfg JoinGrantConfig) (JoinGrantClaims, error) {
	grant = strings.TrimSpace(grant)
	if grant == "" {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant is required")
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Issuer == "" || cfg.Audience == "" || len(cfg.Key) != ed25519.PublicKeySize {
		return JoinGrantClaims{}, errors.New("join grant verifier is not configured")
	}

	parsed, err := parseJoinGrant(grant)
	if err != nil {
		return JoinGrantClaims{}, err
	}
	if parsed.Header.Alg != "EdDSA" {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant alg is invalid")
	}
	if !ed25519.Verify(cfg.Key, []byte(parsed.SigningInput), parsed.Signature) {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant signature is invalid")
	}

	if parsed.Payload.Issuer == "" || parsed.Payload.Issuer != cfg.Issuer {
		return JoinGrantClaims{}, apperrors.WithMetadata(
			apperrors.CodeInviteJoinGrantMismatch,
			"join grant issuer mismatch",
			map[string]string{"Field": "issuer"},
		)
	}
	if !parsed.Payload.Audience.Contains(cfg.Audience) {
		return JoinGrantClaims{}, apperrors.WithMetadata(
			apperrors.CodeInviteJoinGrantMismatch,
			"join grant audience mismatch",
			map[string]string{"Field": "audience"},
		)
	}

	if parsed.Payload.JWTID == "" {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant jti is required")
	}
	if parsed.Payload.ExpiresAt == 0 {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant exp is required")
	}

	now := cfg.Now().UTC()
	exp := time.Unix(parsed.Payload.ExpiresAt, 0).UTC()
	if !exp.After(now) {
		return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantExpired, "join grant is expired")
	}
	if parsed.Payload.NotBefore > 0 {
		nbf := time.Unix(parsed.Payload.NotBefore, 0).UTC()
		if now.Before(nbf) {
			return JoinGrantClaims{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant not active yet")
		}
	}

	if strings.TrimSpace(parsed.Payload.CampaignID) == "" || parsed.Payload.CampaignID != expected.CampaignID {
		return JoinGrantClaims{}, apperrors.WithMetadata(
			apperrors.CodeInviteJoinGrantMismatch,
			"join grant campaign mismatch",
			map[string]string{"Field": "campaign_id"},
		)
	}
	if strings.TrimSpace(parsed.Payload.InviteID) == "" || parsed.Payload.InviteID != expected.InviteID {
		return JoinGrantClaims{}, apperrors.WithMetadata(
			apperrors.CodeInviteJoinGrantMismatch,
			"join grant invite mismatch",
			map[string]string{"Field": "invite_id"},
		)
	}
	if strings.TrimSpace(parsed.Payload.UserID) == "" || parsed.Payload.UserID != expected.UserID {
		return JoinGrantClaims{}, apperrors.WithMetadata(
			apperrors.CodeInviteJoinGrantMismatch,
			"join grant user mismatch",
			map[string]string{"Field": "user_id"},
		)
	}

	claims := JoinGrantClaims{
		Issuer:     parsed.Payload.Issuer,
		Audience:   parsed.Payload.Audience.Values,
		ExpiresAt:  exp,
		JWTID:      parsed.Payload.JWTID,
		CampaignID: parsed.Payload.CampaignID,
		InviteID:   parsed.Payload.InviteID,
		UserID:     parsed.Payload.UserID,
	}
	if parsed.Payload.NotBefore > 0 {
		claims.NotBefore = time.Unix(parsed.Payload.NotBefore, 0).UTC()
	}
	if parsed.Payload.IssuedAt > 0 {
		claims.IssuedAt = time.Unix(parsed.Payload.IssuedAt, 0).UTC()
	}
	return claims, nil
}

type joinGrantHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ,omitempty"`
}

type joinGrantPayload struct {
	Issuer     string        `json:"iss"`
	Audience   audienceClaim `json:"aud"`
	ExpiresAt  int64         `json:"exp"`
	NotBefore  int64         `json:"nbf,omitempty"`
	IssuedAt   int64         `json:"iat,omitempty"`
	JWTID      string        `json:"jti,omitempty"`
	CampaignID string        `json:"campaign_id"`
	InviteID   string        `json:"invite_id"`
	UserID     string        `json:"user_id"`
}

type joinGrantToken struct {
	Header       joinGrantHeader
	Payload      joinGrantPayload
	Signature    []byte
	SigningInput string
}

type audienceClaim struct {
	Values []string
}

func (a *audienceClaim) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		value = strings.TrimSpace(value)
		if value != "" {
			a.Values = []string{value}
		}
		return nil
	}
	if data[0] == '[' {
		var values []string
		if err := json.Unmarshal(data, &values); err != nil {
			return err
		}
		for _, value := range values {
			trimmed := strings.TrimSpace(value)
			if trimmed != "" {
				a.Values = append(a.Values, trimmed)
			}
		}
		return nil
	}
	return fmt.Errorf("aud must be string or array")
}

func (a audienceClaim) Contains(value string) bool {
	for _, item := range a.Values {
		if item == value {
			return true
		}
	}
	return false
}

func parseJoinGrant(grant string) (joinGrantToken, error) {
	parts := strings.Split(grant, ".")
	if len(parts) != 3 {
		return joinGrantToken{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant must have three segments")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return joinGrantToken{}, apperrors.Wrap(apperrors.CodeInviteJoinGrantInvalid, "decode join grant header", err)
	}
	var header joinGrantHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return joinGrantToken{}, apperrors.Wrap(apperrors.CodeInviteJoinGrantInvalid, "decode join grant header", err)
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return joinGrantToken{}, apperrors.Wrap(apperrors.CodeInviteJoinGrantInvalid, "decode join grant payload", err)
	}
	var payload joinGrantPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return joinGrantToken{}, apperrors.Wrap(apperrors.CodeInviteJoinGrantInvalid, "decode join grant payload", err)
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return joinGrantToken{}, apperrors.Wrap(apperrors.CodeInviteJoinGrantInvalid, "decode join grant signature", err)
	}
	if len(sig) != ed25519.SignatureSize {
		return joinGrantToken{}, apperrors.New(apperrors.CodeInviteJoinGrantInvalid, "join grant signature size is invalid")
	}

	return joinGrantToken{
		Header:       header,
		Payload:      payload,
		Signature:    sig,
		SigningInput: parts[0] + "." + parts[1],
	}, nil
}

func decodeBase64(value string) ([]byte, error) {
	if value == "" {
		return nil, errors.New("empty base64 value")
	}
	decoded, err := base64.RawStdEncoding.DecodeString(value)
	if err == nil {
		return decoded, nil
	}
	return base64.StdEncoding.DecodeString(value)
}
