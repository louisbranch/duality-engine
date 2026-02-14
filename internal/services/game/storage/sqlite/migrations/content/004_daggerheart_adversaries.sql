-- +migrate Up

CREATE TABLE daggerheart_adversary_entries (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tier INTEGER NOT NULL DEFAULT 0,
    role TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    motives TEXT NOT NULL DEFAULT '',
    difficulty INTEGER NOT NULL DEFAULT 0,
    major_threshold INTEGER NOT NULL DEFAULT 0,
    severe_threshold INTEGER NOT NULL DEFAULT 0,
    hp INTEGER NOT NULL DEFAULT 0,
    stress INTEGER NOT NULL DEFAULT 0,
    armor INTEGER NOT NULL DEFAULT 0,
    attack_modifier INTEGER NOT NULL DEFAULT 0,
    standard_attack_json TEXT NOT NULL DEFAULT '{}',
    experiences_json TEXT NOT NULL DEFAULT '[]',
    features_json TEXT NOT NULL DEFAULT '[]',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daggerheart_adversary_entries;
