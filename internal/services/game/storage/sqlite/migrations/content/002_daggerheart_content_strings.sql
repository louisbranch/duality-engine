-- +migrate Up

CREATE TABLE daggerheart_content_strings (
    content_id TEXT NOT NULL,
    content_type TEXT NOT NULL,
    field TEXT NOT NULL,
    locale TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    PRIMARY KEY (content_id, field, locale)
);

-- +migrate Down

DROP TABLE IF EXISTS daggerheart_content_strings;
