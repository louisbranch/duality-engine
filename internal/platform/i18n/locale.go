// Package i18n provides locale helpers for the platform.
package i18n

import (
	"strings"

	commonv1 "github.com/louisbranch/fracturing.space/api/gen/go/common/v1"
	"golang.org/x/text/language"
)

var supportedLocales = []commonv1.Locale{
	commonv1.Locale_LOCALE_EN_US,
	commonv1.Locale_LOCALE_PT_BR,
}

var supportedTags = []language.Tag{
	language.MustParse("en-US"),
	language.MustParse("pt-BR"),
}

var tagMatcher = language.NewMatcher(supportedTags)

var localeByTag = map[string]commonv1.Locale{
	"en-US": commonv1.Locale_LOCALE_EN_US,
	"pt-BR": commonv1.Locale_LOCALE_PT_BR,
}

var tagByLocale = map[commonv1.Locale]language.Tag{
	commonv1.Locale_LOCALE_EN_US: language.MustParse("en-US"),
	commonv1.Locale_LOCALE_PT_BR: language.MustParse("pt-BR"),
}

var localeAliases = map[string]commonv1.Locale{
	"en":    commonv1.Locale_LOCALE_EN_US,
	"en-US": commonv1.Locale_LOCALE_EN_US,
	"pt":    commonv1.Locale_LOCALE_PT_BR,
	"pt-BR": commonv1.Locale_LOCALE_PT_BR,
}

// SupportedLocales returns the supported locale enums.
func SupportedLocales() []commonv1.Locale {
	locales := make([]commonv1.Locale, len(supportedLocales))
	copy(locales, supportedLocales)
	return locales
}

// SupportedTags returns the supported locale tags.
func SupportedTags() []language.Tag {
	tags := make([]language.Tag, len(supportedTags))
	copy(tags, supportedTags)
	return tags
}

// DefaultLocale returns the default locale when one is not provided.
func DefaultLocale() commonv1.Locale {
	return commonv1.Locale_LOCALE_EN_US
}

// NormalizeLocale coerces unknown locales to the default.
func NormalizeLocale(locale commonv1.Locale) commonv1.Locale {
	if _, ok := tagByLocale[locale]; ok {
		return locale
	}
	return DefaultLocale()
}

// LocaleString returns the BCP-47 string for the locale.
func LocaleString(locale commonv1.Locale) string {
	return TagForLocale(locale).String()
}

// DefaultTag returns the default locale as a language tag.
func DefaultTag() language.Tag {
	return TagForLocale(DefaultLocale())
}

// TagForLocale returns the language tag for a locale.
func TagForLocale(locale commonv1.Locale) language.Tag {
	normalized := NormalizeLocale(locale)
	return tagByLocale[normalized]
}

// LocaleForTag returns the locale enum for a language tag.
func LocaleForTag(tag language.Tag) commonv1.Locale {
	if locale, ok := localeByTag[tag.String()]; ok {
		return locale
	}
	return DefaultLocale()
}

// ParseLocale parses an explicit locale value into the enum.
func ParseLocale(value string) (commonv1.Locale, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return DefaultLocale(), false
	}
	if locale, ok := localeAliases[trimmed]; ok {
		return locale, true
	}
	parsed, err := language.Parse(trimmed)
	if err != nil {
		return DefaultLocale(), false
	}
	if locale, ok := localeAliases[parsed.String()]; ok {
		return locale, true
	}
	return DefaultLocale(), false
}

// ParseTag parses an explicit locale value into a language tag.
func ParseTag(value string) (language.Tag, bool) {
	locale, ok := ParseLocale(value)
	if !ok {
		return language.Tag{}, false
	}
	return TagForLocale(locale), true
}

// MatchTags returns the best supported tag for the provided list.
func MatchTags(tags []language.Tag) language.Tag {
	if len(tags) == 0 {
		return DefaultTag()
	}
	matched, _, _ := tagMatcher.Match(tags...)
	return matched
}
