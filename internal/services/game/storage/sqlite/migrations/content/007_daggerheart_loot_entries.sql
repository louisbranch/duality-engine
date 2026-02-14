-- +migrate Up

CREATE TABLE daggerheart_loot_entries (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    roll INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daggerheart_loot_entries;
