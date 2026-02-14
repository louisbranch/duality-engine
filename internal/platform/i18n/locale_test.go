package i18n

import (
	"testing"

	commonv1 "github.com/louisbranch/fracturing.space/api/gen/go/common/v1"
	"golang.org/x/text/language"
)

func TestSupportedLocales(t *testing.T) {
	locales := SupportedLocales()
	if len(locales) == 0 {
		t.Fatal("expected at least one supported locale")
	}
	// Modifying the returned slice should not affect the source.
	original := len(locales)
	locales[0] = commonv1.Locale_LOCALE_UNSPECIFIED
	if got := SupportedLocales(); len(got) != original {
		t.Fatalf("mutation leaked: expected %d locales, got %d", original, len(got))
	}
}

func TestSupportedTags(t *testing.T) {
	tags := SupportedTags()
	if len(tags) == 0 {
		t.Fatal("expected at least one supported tag")
	}
}

func TestDefaultLocaleAndTag(t *testing.T) {
	if DefaultLocale() == commonv1.Locale_LOCALE_UNSPECIFIED {
		t.Fatal("expected a specific default locale")
	}
	if DefaultTag() == (language.Tag{}) {
		t.Fatal("expected a valid default tag")
	}
}

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		name  string
		input commonv1.Locale
		want  commonv1.Locale
	}{
		{name: "supported", input: commonv1.Locale_LOCALE_EN_US, want: commonv1.Locale_LOCALE_EN_US},
		{name: "unsupported", input: commonv1.Locale_LOCALE_UNSPECIFIED, want: DefaultLocale()},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := NormalizeLocale(test.input)
			if got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestLocaleTagMappings(t *testing.T) {
	tag := TagForLocale(commonv1.Locale_LOCALE_PT_BR)
	if tag.String() != "pt-BR" {
		t.Fatalf("expected tag pt-BR, got %q", tag.String())
	}

	locale := LocaleForTag(language.MustParse("en-US"))
	if locale != commonv1.Locale_LOCALE_EN_US {
		t.Fatalf("expected locale %v, got %v", commonv1.Locale_LOCALE_EN_US, locale)
	}

	fallback := LocaleForTag(language.MustParse("fr-FR"))
	if fallback != DefaultLocale() {
		t.Fatalf("expected default locale %v, got %v", DefaultLocale(), fallback)
	}
}

func TestLocaleString(t *testing.T) {
	got := LocaleString(commonv1.Locale_LOCALE_EN_US)
	if got != "en-US" {
		t.Fatalf("expected en-US, got %q", got)
	}
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  commonv1.Locale
		ok    bool
	}{
		{name: "empty", input: "", want: DefaultLocale(), ok: false},
		{name: "alias-en", input: "en", want: commonv1.Locale_LOCALE_EN_US, ok: true},
		{name: "alias-pt", input: "pt", want: commonv1.Locale_LOCALE_PT_BR, ok: true},
		{name: "canonical", input: "pt-BR", want: commonv1.Locale_LOCALE_PT_BR, ok: true},
		{name: "parsed", input: "en-us", want: commonv1.Locale_LOCALE_EN_US, ok: true},
		{name: "invalid", input: "fr", want: DefaultLocale(), ok: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, ok := ParseLocale(test.input)
			if ok != test.ok {
				t.Fatalf("expected ok %v, got %v", test.ok, ok)
			}
			if got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestParseTag(t *testing.T) {
	tag, ok := ParseTag("pt")
	if !ok {
		t.Fatal("expected ok")
	}
	if tag.String() != "pt-BR" {
		t.Fatalf("expected pt-BR, got %q", tag.String())
	}

	_, ok = ParseTag("fr")
	if ok {
		t.Fatal("expected not ok")
	}
}

func TestMatchTags(t *testing.T) {
	assertBaseRegion := func(t *testing.T, tag language.Tag, wantBase, wantRegion string) {
		t.Helper()
		base, _ := tag.Base()
		if base.String() != wantBase {
			t.Fatalf("expected base %q, got %q", wantBase, base.String())
		}
		region, _ := tag.Region()
		if region.String() != wantRegion {
			t.Fatalf("expected region %q, got %q", wantRegion, region.String())
		}
	}

	tests := []struct {
		name       string
		tags       []language.Tag
		wantBase   string
		wantRegion string
		wantTag    language.Tag
	}{
		{
			name:    "empty",
			tags:    nil,
			wantTag: DefaultTag(),
		},
		{
			name:       "english",
			tags:       []language.Tag{language.MustParse("en-GB")},
			wantBase:   "en",
			wantRegion: "US",
		},
		{
			name:       "portuguese",
			tags:       []language.Tag{language.MustParse("pt-PT"), language.MustParse("en-US")},
			wantBase:   "pt",
			wantRegion: "BR",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := MatchTags(test.tags)
			if test.wantTag != (language.Tag{}) {
				if got != test.wantTag {
					t.Fatalf("expected %q, got %q", test.wantTag, got)
				}
				return
			}
			assertBaseRegion(t, got, test.wantBase, test.wantRegion)
		})
	}
}
