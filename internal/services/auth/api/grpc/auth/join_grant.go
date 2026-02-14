package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	envJoinGrantIssuer     = "FRACTURING_SPACE_JOIN_GRANT_ISSUER"
	envJoinGrantAudience   = "FRACTURING_SPACE_JOIN_GRANT_AUDIENCE"
	envJoinGrantPrivateKey = "FRACTURING_SPACE_JOIN_GRANT_PRIVATE_KEY"
	envJoinGrantTTL        = "FRACTURING_SPACE_JOIN_GRANT_TTL"
)

const defaultJoinGrantTTL = 5 * time.Minute

type joinGrantConfig struct {
	issuer   string
	audience string
	key      ed25519.PrivateKey
	ttl      time.Duration
}

func loadJoinGrantConfigFromEnv() (joinGrantConfig, error) {
	issuer := strings.TrimSpace(os.Getenv(envJoinGrantIssuer))
	if issuer == "" {
		return joinGrantConfig{}, fmt.Errorf("%s is required", envJoinGrantIssuer)
	}
	audience := strings.TrimSpace(os.Getenv(envJoinGrantAudience))
	if audience == "" {
		return joinGrantConfig{}, fmt.Errorf("%s is required", envJoinGrantAudience)
	}
	keyRaw := strings.TrimSpace(os.Getenv(envJoinGrantPrivateKey))
	if keyRaw == "" {
		return joinGrantConfig{}, fmt.Errorf("%s is required", envJoinGrantPrivateKey)
	}
	keyBytes, err := decodeBase64(keyRaw)
	if err != nil {
		return joinGrantConfig{}, fmt.Errorf("decode join grant private key: %w", err)
	}
	if len(keyBytes) != ed25519.PrivateKeySize {
		return joinGrantConfig{}, fmt.Errorf("join grant private key must be %d bytes", ed25519.PrivateKeySize)
	}

	ttl := defaultJoinGrantTTL
	if rawTTL := strings.TrimSpace(os.Getenv(envJoinGrantTTL)); rawTTL != "" {
		parsed, err := time.ParseDuration(rawTTL)
		if err != nil {
			return joinGrantConfig{}, fmt.Errorf("parse join grant ttl: %w", err)
		}
		if parsed <= 0 {
			return joinGrantConfig{}, fmt.Errorf("join grant ttl must be positive")
		}
		ttl = parsed
	}

	return joinGrantConfig{
		issuer:   issuer,
		audience: audience,
		key:      ed25519.PrivateKey(keyBytes),
		ttl:      ttl,
	}, nil
}

func encodeJoinGrant(cfg joinGrantConfig, payload map[string]any) (string, error) {
	if cfg.issuer == "" || cfg.audience == "" || len(cfg.key) != ed25519.PrivateKeySize {
		return "", errors.New("join grant signer is not configured")
	}
	headerJSON, err := json.Marshal(map[string]string{
		"alg": "EdDSA",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("encode join grant header: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode join grant payload: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := encodedHeader + "." + encodedPayload
	signature := ed25519.Sign(cfg.key, []byte(signingInput))
	encodedSig := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + encodedSig, nil
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
