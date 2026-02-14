# Locale Helper Tests

## Coverage
- Unit tests cover locale parsing, normalization, tag matching, and mappings in `internal/platform/i18n/locale_test.go`.
- Match behavior asserts base/region to allow for Unicode extension additions from the matcher.

## Integration Note
The integration target runs `event-catalog-check`, which regenerates `docs/events/event-catalog.md` and fails if the file differs from the index. After `go generate ./internal/services/game/domain/campaign/event`, stage the catalog file before running `make integration`.
