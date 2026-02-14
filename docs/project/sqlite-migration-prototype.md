## SQLite migrations in prototype mode

This project is in prototype mode, so SQLite migrations are intentionally non-incremental.

- Prefer single, full-schema migrations per store.
- Avoid ALTER-based migrations; drop and recreate tables instead.
- Keep migration order simple and reasoning-friendly by collapsing superseded files.
