-- +migrate Up

CREATE TABLE daggerheart_beastforms (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tier INTEGER NOT NULL DEFAULT 0,
    examples TEXT NOT NULL DEFAULT '',
    trait TEXT NOT NULL DEFAULT '',
    trait_bonus INTEGER NOT NULL DEFAULT 0,
    evasion_bonus INTEGER NOT NULL DEFAULT 0,
    attack_json TEXT NOT NULL DEFAULT '{}',
    advantages_json TEXT NOT NULL DEFAULT '[]',
    features_json TEXT NOT NULL DEFAULT '[]',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daggerheart_beastforms;
