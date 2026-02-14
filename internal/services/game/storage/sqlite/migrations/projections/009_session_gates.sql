-- +migrate Up

CREATE TABLE session_gates (
    campaign_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    gate_id TEXT NOT NULL,
    gate_type TEXT NOT NULL,
    status TEXT NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL,
    created_by_actor_type TEXT NOT NULL,
    created_by_actor_id TEXT NOT NULL DEFAULT '',
    resolved_at INTEGER,
    resolved_by_actor_type TEXT,
    resolved_by_actor_id TEXT,
    metadata_json BLOB,
    resolution_json BLOB,
    PRIMARY KEY (campaign_id, session_id, gate_id)
);

CREATE INDEX idx_session_gates_open ON session_gates(campaign_id, session_id, status);

-- +migrate Down

DROP INDEX IF EXISTS idx_session_gates_open;
DROP TABLE IF EXISTS session_gates;
