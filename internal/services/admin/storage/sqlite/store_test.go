package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOpenRequiresPath(t *testing.T) {
	if _, err := Open(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestPutUserSessionStoresTimestamp(t *testing.T) {
	store := openTempStore(t)

	createdAt := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	if err := store.PutUserSession(context.Background(), "session-1", createdAt); err != nil {
		t.Fatalf("put user session: %v", err)
	}

	var storedID string
	var storedAt string
	row := store.sqlDB.QueryRow("SELECT session_id, created_at FROM user_sessions WHERE session_id = ?", "session-1")
	if err := row.Scan(&storedID, &storedAt); err != nil {
		t.Fatalf("scan user session: %v", err)
	}
	if storedID != "session-1" {
		t.Fatalf("expected session id session-1, got %s", storedID)
	}
	if storedAt != createdAt.Format(timeFormat) {
		t.Fatalf("expected created_at %s, got %s", createdAt.Format(timeFormat), storedAt)
	}
}

func TestPutUserSessionDefaultsTime(t *testing.T) {
	store := openTempStore(t)

	if err := store.PutUserSession(context.Background(), "session-2", time.Time{}); err != nil {
		t.Fatalf("put user session: %v", err)
	}

	var storedAt string
	row := store.sqlDB.QueryRow("SELECT created_at FROM user_sessions WHERE session_id = ?", "session-2")
	if err := row.Scan(&storedAt); err != nil {
		t.Fatalf("scan user session: %v", err)
	}
	if storedAt == "" {
		t.Fatal("expected created_at to be set")
	}
}

func TestPutUserSessionValidation(t *testing.T) {
	store := openTempStore(t)

	if err := store.PutUserSession(context.Background(), "", time.Now()); err == nil {
		t.Fatal("expected error for empty session id")
	}
}

func TestPutUserSessionRequiresStore(t *testing.T) {
	var store *Store
	if err := store.PutUserSession(context.Background(), "session-3", time.Now()); err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestCloseNilSafe(t *testing.T) {
	var s *Store
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error from nil store: %v", err)
	}

	s = &Store{}
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected error from store with nil db: %v", err)
	}
}

func TestPutUserSessionContextCancelled(t *testing.T) {
	store := openTempStore(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.PutUserSession(ctx, "session-1", time.Now()); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestExtractUpMigration(t *testing.T) {
	t.Run("no markers", func(t *testing.T) {
		content := "CREATE TABLE foo (id TEXT);"
		up := extractUpMigration(content)
		if up != content {
			t.Errorf("expected full content, got %q", up)
		}
	})

	t.Run("up only", func(t *testing.T) {
		content := "-- +migrate Up\nCREATE TABLE foo (id TEXT);"
		up := extractUpMigration(content)
		if !strings.Contains(up, "CREATE TABLE") {
			t.Errorf("expected CREATE TABLE, got %q", up)
		}
	})

	t.Run("up and down", func(t *testing.T) {
		content := "-- +migrate Up\nCREATE TABLE foo (id TEXT);\n-- +migrate Down\nDROP TABLE foo;"
		up := extractUpMigration(content)
		if !strings.Contains(up, "CREATE TABLE") {
			t.Errorf("expected CREATE TABLE, got %q", up)
		}
		if strings.Contains(up, "DROP TABLE") {
			t.Error("did not expect DROP TABLE in up migration")
		}
	})
}

func TestIsAlreadyExistsError(t *testing.T) {
	if !isAlreadyExistsError(fmt.Errorf("table already exists")) {
		t.Error("expected true for already exists error")
	}
	if isAlreadyExistsError(fmt.Errorf("not found")) {
		t.Error("expected false for unrelated error")
	}
}

func openTempStore(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "admin.db")
	store, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil && err != sql.ErrConnDone {
			t.Fatalf("close store: %v", err)
		}
	})
	return store
}
