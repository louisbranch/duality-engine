-- +migrate Up

CREATE TABLE session_spotlight (
    campaign_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    spotlight_type TEXT NOT NULL,
    character_id TEXT NOT NULL DEFAULT '',
    updated_at INTEGER NOT NULL,
    updated_by_actor_type TEXT NOT NULL,
    updated_by_actor_id TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (campaign_id, session_id)
);

CREATE INDEX idx_session_spotlight_session ON session_spotlight(campaign_id, session_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_session_spotlight_session;
DROP TABLE IF EXISTS session_spotlight;
