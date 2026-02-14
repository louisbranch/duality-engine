package maintenance

import (
	"bytes"
	"encoding/json"
	"flag"
	"reflect"
	"strings"
	"testing"
)

func TestResolveCampaignIDs(t *testing.T) {
	tests := []struct {
		single   string
		list     string
		expected []string
		wantErr  bool
	}{
		{single: "", list: "", wantErr: true},
		{single: "c1", list: "c2", wantErr: true},
		{single: "c1", list: "", expected: []string{"c1"}},
		{single: "", list: "c1, c2", expected: []string{"c1", "c2"}},
		{single: "", list: " , c1 , , c2 ", expected: []string{"c1", "c2"}},
	}

	for _, tc := range tests {
		got, err := resolveCampaignIDs(tc.single, tc.list)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected error for %q/%q", tc.single, tc.list)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q/%q: %v", tc.single, tc.list, err)
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, got)
		}
	}
}

func TestSplitCSV(t *testing.T) {
	if got := splitCSV(" a, b ,, "); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("expected trimmed entries, got %v", got)
	}
}

func TestCapWarnings(t *testing.T) {
	warnings := []string{"a", "b", "c"}
	if got, total := capWarnings(warnings, 0); total != 3 || len(got) != 3 {
		t.Fatalf("expected all warnings, got %v (total=%d)", got, total)
	}
	if got, total := capWarnings(warnings, 2); total != 3 || len(got) != 2 {
		t.Fatalf("expected capped warnings, got %v (total=%d)", got, total)
	}
}

func TestParseConfigDefaults(t *testing.T) {
	fs := flag.NewFlagSet("maintenance", flag.ContinueOnError)
	cfg, err := ParseConfig(fs, nil, func(string) (string, bool) { return "", false })
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	if cfg.EventsDBPath != "data/game-events.db" {
		t.Fatalf("expected default events db path, got %q", cfg.EventsDBPath)
	}
	if cfg.ProjectionsDBPath != "data/game-projections.db" {
		t.Fatalf("expected default projections db path, got %q", cfg.ProjectionsDBPath)
	}
	if cfg.WarningsCap != 25 {
		t.Fatalf("expected warnings cap 25, got %d", cfg.WarningsCap)
	}
}

func TestParseConfigOverrides(t *testing.T) {
	fs := flag.NewFlagSet("maintenance", flag.ContinueOnError)
	lookup := func(key string) (string, bool) {
		switch key {
		case "FRACTURING_SPACE_GAME_EVENTS_DB_PATH":
			return "env-events", true
		case "FRACTURING_SPACE_GAME_PROJECTIONS_DB_PATH":
			return "env-projections", true
		default:
			return "", false
		}
	}
	args := []string{
		"-events-db-path", "flag-events",
		"-projections-db-path", "flag-projections",
		"-warnings-cap", "5",
	}
	cfg, err := ParseConfig(fs, args, lookup)
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	if cfg.EventsDBPath != "flag-events" {
		t.Fatalf("expected flag override for events db, got %q", cfg.EventsDBPath)
	}
	if cfg.ProjectionsDBPath != "flag-projections" {
		t.Fatalf("expected flag override for projections db, got %q", cfg.ProjectionsDBPath)
	}
	if cfg.WarningsCap != 5 {
		t.Fatalf("expected warnings cap 5, got %d", cfg.WarningsCap)
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Run("nil lookup returns fallback", func(t *testing.T) {
		got := envOrDefault(nil, []string{"KEY"}, "fb")
		if got != "fb" {
			t.Errorf("expected %q, got %q", "fb", got)
		}
	})

	t.Run("key found", func(t *testing.T) {
		lookup := func(key string) (string, bool) {
			if key == "A" {
				return "found", true
			}
			return "", false
		}
		got := envOrDefault(lookup, []string{"A"}, "fb")
		if got != "found" {
			t.Errorf("expected %q, got %q", "found", got)
		}
	})

	t.Run("whitespace value falls through", func(t *testing.T) {
		lookup := func(key string) (string, bool) {
			if key == "A" {
				return "  ", true
			}
			return "", false
		}
		got := envOrDefault(lookup, []string{"A"}, "fb")
		if got != "fb" {
			t.Errorf("expected %q, got %q", "fb", got)
		}
	})

	t.Run("first matching key wins", func(t *testing.T) {
		lookup := func(key string) (string, bool) {
			switch key {
			case "A":
				return "", false
			case "B":
				return "b-val", true
			default:
				return "", false
			}
		}
		got := envOrDefault(lookup, []string{"A", "B"}, "fb")
		if got != "b-val" {
			t.Errorf("expected %q, got %q", "b-val", got)
		}
	})

	t.Run("no keys returns fallback", func(t *testing.T) {
		lookup := func(string) (string, bool) { return "", false }
		got := envOrDefault(lookup, nil, "fb")
		if got != "fb" {
			t.Errorf("expected %q, got %q", "fb", got)
		}
	})
}

func TestDefaultEventsDBPath(t *testing.T) {
	t.Run("no env", func(t *testing.T) {
		lookup := func(string) (string, bool) { return "", false }
		got := defaultEventsDBPath(lookup)
		if got != "data/game-events.db" {
			t.Errorf("expected default path, got %q", got)
		}
	})

	t.Run("env set", func(t *testing.T) {
		lookup := func(key string) (string, bool) {
			if key == "FRACTURING_SPACE_GAME_EVENTS_DB_PATH" {
				return "/custom/events.db", true
			}
			return "", false
		}
		got := defaultEventsDBPath(lookup)
		if got != "/custom/events.db" {
			t.Errorf("expected %q, got %q", "/custom/events.db", got)
		}
	})
}

func TestDefaultProjectionsDBPath(t *testing.T) {
	t.Run("no env", func(t *testing.T) {
		lookup := func(string) (string, bool) { return "", false }
		got := defaultProjectionsDBPath(lookup)
		if got != "data/game-projections.db" {
			t.Errorf("expected default path, got %q", got)
		}
	})

	t.Run("env set", func(t *testing.T) {
		lookup := func(key string) (string, bool) {
			if key == "FRACTURING_SPACE_GAME_PROJECTIONS_DB_PATH" {
				return "/custom/proj.db", true
			}
			return "", false
		}
		got := defaultProjectionsDBPath(lookup)
		if got != "/custom/proj.db" {
			t.Errorf("expected %q, got %q", "/custom/proj.db", got)
		}
	})
}

func TestOutputJSON(t *testing.T) {
	t.Run("valid result", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{
			CampaignID: "c1",
			Mode:       "scan",
		}
		outputJSON(&out, &errOut, result)
		if errOut.Len() != 0 {
			t.Errorf("unexpected error output: %s", errOut.String())
		}
		var decoded runResult
		if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
			t.Fatalf("invalid JSON output: %v", err)
		}
		if decoded.CampaignID != "c1" {
			t.Errorf("campaign_id = %q, want %q", decoded.CampaignID, "c1")
		}
	})

	t.Run("with warnings", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{
			CampaignID:    "c2",
			Mode:          "validate",
			Warnings:      []string{"warn1"},
			WarningsTotal: 5,
		}
		outputJSON(&out, &errOut, result)
		if !strings.Contains(out.String(), `"warnings_total":5`) {
			t.Errorf("expected warnings_total in output: %s", out.String())
		}
	})
}

func TestPrintResult(t *testing.T) {
	t.Run("error output", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{Error: "something failed"}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(errOut.String(), "Error: something failed") {
			t.Errorf("expected error in errOut: %s", errOut.String())
		}
	})

	t.Run("warnings output", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{Warnings: []string{"w1", "w2"}, WarningsTotal: 5}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(errOut.String(), "Warning: w1") {
			t.Errorf("expected warning w1: %s", errOut.String())
		}
		if !strings.Contains(errOut.String(), "3 more warnings suppressed") {
			t.Errorf("expected suppressed warning count: %s", errOut.String())
		}
	})

	t.Run("integrity report", func(t *testing.T) {
		var out, errOut bytes.Buffer
		report := integrityReport{
			LastSeq:             100,
			CharacterMismatches: 2,
			MissingStates:       1,
			GmFearMatch:         true,
			GmFearSource:        5,
			GmFearReplay:        5,
		}
		reportJSON, _ := json.Marshal(report)
		result := runResult{
			CampaignID: "c1",
			Mode:       "integrity",
			Report:     reportJSON,
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(out.String(), "Integrity check") {
			t.Errorf("expected integrity output: %s", out.String())
		}
		if !strings.Contains(out.String(), "GM fear match: true") {
			t.Errorf("expected GM fear match: %s", out.String())
		}
		if !strings.Contains(out.String(), "Character state mismatches: 2") {
			t.Errorf("expected character mismatches: %s", out.String())
		}
	})

	t.Run("scan report", func(t *testing.T) {
		var out, errOut bytes.Buffer
		report := snapshotScanReport{
			LastSeq:        50,
			TotalEvents:    100,
			SnapshotEvents: 10,
		}
		reportJSON, _ := json.Marshal(report)
		result := runResult{
			CampaignID: "c1",
			Mode:       "scan",
			Report:     reportJSON,
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(out.String(), "Scanned snapshot-related events") {
			t.Errorf("expected scan output: %s", out.String())
		}
	})

	t.Run("validate report", func(t *testing.T) {
		var out, errOut bytes.Buffer
		report := snapshotScanReport{
			LastSeq:        50,
			TotalEvents:    100,
			SnapshotEvents: 10,
			InvalidEvents:  3,
		}
		reportJSON, _ := json.Marshal(report)
		result := runResult{
			CampaignID: "c1",
			Mode:       "validate",
			Report:     reportJSON,
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(out.String(), "Validated snapshot-related events") {
			t.Errorf("expected validate output: %s", out.String())
		}
	})

	t.Run("replay report", func(t *testing.T) {
		var out, errOut bytes.Buffer
		report := snapshotScanReport{LastSeq: 50}
		reportJSON, _ := json.Marshal(report)
		result := runResult{
			CampaignID: "c1",
			Mode:       "replay",
			Report:     reportJSON,
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(out.String(), "Replayed snapshot-related events") {
			t.Errorf("expected replay output: %s", out.String())
		}
	})

	t.Run("prefix applied", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{Error: "oops"}
		printResult(&out, &errOut, result, "[c1] ")
		if !strings.Contains(errOut.String(), "[c1] Error: oops") {
			t.Errorf("expected prefix in output: %s", errOut.String())
		}
	})

	t.Run("empty report returns early", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{CampaignID: "c1", Mode: "scan"}
		printResult(&out, &errOut, result, "")
		if out.Len() != 0 {
			t.Errorf("expected no output for empty report: %s", out.String())
		}
	})

	t.Run("invalid integrity JSON", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{
			CampaignID: "c1",
			Mode:       "integrity",
			Report:     json.RawMessage(`{invalid`),
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(errOut.String(), "decode report") {
			t.Errorf("expected decode error: %s", errOut.String())
		}
	})

	t.Run("invalid scan JSON", func(t *testing.T) {
		var out, errOut bytes.Buffer
		result := runResult{
			CampaignID: "c1",
			Mode:       "scan",
			Report:     json.RawMessage(`{invalid`),
		}
		printResult(&out, &errOut, result, "")
		if !strings.Contains(errOut.String(), "decode report") {
			t.Errorf("expected decode error: %s", errOut.String())
		}
	})
}

func TestRunValidationErrors(t *testing.T) {
	t.Run("integrity with dry-run", func(t *testing.T) {
		cfg := Config{
			CampaignID: "c1",
			Integrity:  true,
			DryRun:     true,
		}
		err := Run(t.Context(), cfg, nil, nil)
		if err == nil || !strings.Contains(err.Error(), "-integrity cannot be combined") {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("integrity with validate", func(t *testing.T) {
		cfg := Config{
			CampaignID: "c1",
			Integrity:  true,
			Validate:   true,
		}
		err := Run(t.Context(), cfg, nil, nil)
		if err == nil || !strings.Contains(err.Error(), "-integrity cannot be combined") {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("integrity with after-seq", func(t *testing.T) {
		cfg := Config{
			CampaignID: "c1",
			Integrity:  true,
			AfterSeq:   10,
		}
		err := Run(t.Context(), cfg, nil, nil)
		if err == nil || !strings.Contains(err.Error(), "-integrity does not support -after-seq") {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("no campaign IDs", func(t *testing.T) {
		cfg := Config{}
		err := Run(t.Context(), cfg, nil, nil)
		if err == nil {
			t.Fatal("expected error for no campaign IDs")
		}
	})

	t.Run("negative warnings cap", func(t *testing.T) {
		cfg := Config{
			CampaignID:  "c1",
			WarningsCap: -1,
		}
		err := Run(t.Context(), cfg, nil, nil)
		if err == nil || !strings.Contains(err.Error(), "warnings-cap") {
			t.Fatalf("expected warnings-cap error, got %v", err)
		}
	})
}
