-- +migrate Up

CREATE TABLE daggerheart_companion_experiences (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daggerheart_companion_experiences;
