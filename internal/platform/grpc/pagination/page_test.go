package pagination

import "testing"

func TestClampPageSize(t *testing.T) {
	tests := []struct {
		name  string
		value int32
		cfg   PageSizeConfig
		want  int
	}{
		{"positive within bounds", 10, PageSizeConfig{Default: 20, Max: 50}, 10},
		{"zero uses default", 0, PageSizeConfig{Default: 20, Max: 50}, 20},
		{"negative uses default", -5, PageSizeConfig{Default: 20, Max: 50}, 20},
		{"exceeds max", 100, PageSizeConfig{Default: 20, Max: 50}, 50},
		{"equal to max", 50, PageSizeConfig{Default: 20, Max: 50}, 50},
		{"no max limit", 999, PageSizeConfig{Default: 20, Max: 0}, 999},
		{"zero default and zero value falls back to 1", 0, PageSizeConfig{Default: 0, Max: 50}, 1},
		{"zero default and negative value falls back to 1", -1, PageSizeConfig{Default: 0, Max: 0}, 1},
		{"one", 1, PageSizeConfig{Default: 20, Max: 50}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ClampPageSize(tc.value, tc.cfg)
			if got != tc.want {
				t.Errorf("ClampPageSize(%d, %+v) = %d, want %d", tc.value, tc.cfg, got, tc.want)
			}
		})
	}
}

func TestNormalizeOrderBy(t *testing.T) {
	cfg := OrderByConfig{
		Default: "created_at DESC",
		Allowed: []string{"created_at ASC", "created_at DESC", "name ASC"},
	}

	t.Run("empty uses default", func(t *testing.T) {
		got, err := NormalizeOrderBy("", cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got != "created_at DESC" {
			t.Errorf("got %q, want %q", got, "created_at DESC")
		}
	})

	t.Run("allowed value", func(t *testing.T) {
		got, err := NormalizeOrderBy("name ASC", cfg)
		if err != nil {
			t.Fatal(err)
		}
		if got != "name ASC" {
			t.Errorf("got %q, want %q", got, "name ASC")
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := NormalizeOrderBy("invalid", cfg)
		if err == nil {
			t.Error("expected error for invalid order_by")
		}
	})

	t.Run("empty allowed list rejects non-empty", func(t *testing.T) {
		_, err := NormalizeOrderBy("name ASC", OrderByConfig{Default: "id"})
		if err == nil {
			t.Error("expected error for unrecognized order_by with empty allowed list")
		}
	})

	t.Run("empty allowed list allows empty", func(t *testing.T) {
		got, err := NormalizeOrderBy("", OrderByConfig{Default: "id"})
		if err != nil {
			t.Fatal(err)
		}
		if got != "id" {
			t.Errorf("got %q, want %q", got, "id")
		}
	})
}
