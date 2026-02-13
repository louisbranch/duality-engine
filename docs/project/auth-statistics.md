# Auth Statistics

The auth service exposes aggregate counts for auth users via `auth.v1.StatisticsService`.

## Scope

- Aggregates are computed over the auth user set.
- `since` filters by user `created_at`; omitted means all time.

## Storage

- SQLite computes the count directly from the `users` table.
- Auth statistics are read-only and do not introduce new schema.
