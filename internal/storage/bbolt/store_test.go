package bbolt

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/louisbranch/duality-engine/internal/campaign/domain"
	"github.com/louisbranch/duality-engine/internal/storage"
)

func TestCampaignStorePutGet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "duality.db")
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()

	now := time.Date(2026, 1, 23, 12, 0, 0, 0, time.UTC)
	campaign := domain.Campaign{
		ID:          "camp-123",
		Name:        "Snowbound",
		GmMode:      domain.GmModeAI,
		PlayerSlots: 4,
		ThemePrompt: "ice and steel",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := store.Put(context.Background(), campaign); err != nil {
		t.Fatalf("put campaign: %v", err)
	}

	loaded, err := store.Get(context.Background(), "camp-123")
	if err != nil {
		t.Fatalf("get campaign: %v", err)
	}
	if loaded.Name != campaign.Name {
		t.Fatalf("expected name %q, got %q", campaign.Name, loaded.Name)
	}
	if loaded.GmMode != campaign.GmMode {
		t.Fatalf("expected gm mode %v, got %v", campaign.GmMode, loaded.GmMode)
	}
	if !loaded.CreatedAt.Equal(now) {
		t.Fatalf("expected created_at %v, got %v", now, loaded.CreatedAt)
	}
}

func TestCampaignStoreGetNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "duality.db")
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()

	_, err = store.Get(context.Background(), "missing")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
