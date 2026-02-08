package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/louisbranch/fracturing.space/internal/web/i18n"
)

// TestWebPageRendering verifies layout rendering based on HTMX requests.
func TestWebPageRendering(t *testing.T) {
	handler := NewHandler(nil)

	tests := []struct {
		name        string
		path        string
		htmx        bool
		contains    []string
		notContains []string
	}{
		{
			name: "home full page",
			path: "/",
			contains: []string{
				"<!doctype html>",
				"Fracturing.Space",
			},
			notContains: []string{
				"<h2>Campaigns</h2>",
			},
		},
		{
			name: "campaigns full page",
			path: "/campaigns",
			contains: []string{
				"<!doctype html>",
				"Fracturing.Space",
				"<h2>Campaigns</h2>",
			},
		},
		{
			name: "campaigns htmx",
			path: "/campaigns",
			htmx: true,
			contains: []string{
				"<h2>Campaigns</h2>",
			},
			notContains: []string{
				"<!doctype html>",
				"Fracturing.Space",
				"<html",
			},
		},
		{
			name: "campaign detail full page",
			path: "/campaigns/camp-123",
			contains: []string{
				"<!doctype html>",
				"Fracturing.Space",
				"Campaign service unavailable.",
				"<h2>Campaign</h2>",
			},
		},
		{
			name: "campaign detail htmx",
			path: "/campaigns/camp-123",
			htmx: true,
			contains: []string{
				"Campaign service unavailable.",
				"<h2>Campaign</h2>",
			},
			notContains: []string{
				"<!doctype html>",
				"Fracturing.Space",
				"<html",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com"+tc.path, nil)
			if tc.htmx {
				req.Header.Set("HX-Request", "true")
			}
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
			}

			body := recorder.Body.String()
			for _, expected := range tc.contains {
				assertContains(t, body, expected)
			}
			for _, unexpected := range tc.notContains {
				assertNotContains(t, body, unexpected)
			}
		})
	}
}

// TestCampaignSessionsRoute verifies session routes render pages correctly.
func TestCampaignSessionsRoute(t *testing.T) {
	handler := NewHandler(nil)

	t.Run("sessions htmx", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/campaigns/camp-123/sessions", nil)
		req.Header.Set("HX-Request", "true")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}

		body := recorder.Body.String()
		assertContains(t, body, "<h3>Sessions</h3>")
		assertNotContains(t, body, "<!doctype html>")
	})

	t.Run("sessions full page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/campaigns/camp-123/sessions", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}

		body := recorder.Body.String()
		assertContains(t, body, "<!doctype html>")
		assertContains(t, body, "<h3>Sessions</h3>")
	})

	t.Run("sessions table htmx", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/campaigns/camp-123/sessions/table", nil)
		req.Header.Set("HX-Request", "true")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}

		body := recorder.Body.String()
		assertContains(t, body, "Session service unavailable.")
	})
}

// assertContains fails the test when the body lacks the expected fragment.
func assertContains(t *testing.T, body string, expected string) {
	t.Helper()
	if !strings.Contains(body, expected) {
		t.Fatalf("expected response to contain %q", expected)
	}
}

// assertNotContains fails the test when the body includes an unexpected fragment.
func assertNotContains(t *testing.T, body string, unexpected string) {
	t.Helper()
	if strings.Contains(body, unexpected) {
		t.Fatalf("expected response to not contain %q", unexpected)
	}
}

// TestEscapeAIP160StringLiteral verifies special character escaping.
func TestEscapeAIP160StringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{`with"quote`, `with\"quote`},
		{`with\backslash`, `with\\backslash`},
		{`both\"chars`, `both\\\"chars`},
		{`a"b\c"d`, `a\"b\\c\"d`},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := escapeAIP160StringLiteral(tc.input)
			if result != tc.expected {
				t.Errorf("escapeAIP160StringLiteral(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestFormatActorType(t *testing.T) {
	loc := i18n.Printer(i18n.Default())

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "system",
			input:    "system",
			expected: loc.Sprintf("filter.actor.system"),
		},
		{
			name:     "participant",
			input:    "participant",
			expected: loc.Sprintf("filter.actor.participant"),
		},
		{
			name:     "gm",
			input:    "gm",
			expected: loc.Sprintf("filter.actor.gm"),
		},
		{
			name:     "fallback",
			input:    "custom",
			expected: "custom",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := formatActorType(tc.input, loc)
			if result != tc.expected {
				t.Errorf("formatActorType(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
