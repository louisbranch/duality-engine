-- +migrate Up

DROP INDEX IF EXISTS idx_participant_claims_participant;
DROP TABLE IF EXISTS participant_claims;

CREATE TABLE participant_claims (
    campaign_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    participant_id TEXT NOT NULL,
    claimed_at INTEGER NOT NULL,
    PRIMARY KEY (campaign_id, user_id),
    UNIQUE (campaign_id, participant_id),
    FOREIGN KEY (campaign_id, participant_id) REFERENCES participants(campaign_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_participant_claims_participant ON participant_claims(participant_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_participant_claims_participant;
DROP TABLE IF EXISTS participant_claims;

CREATE TABLE participant_claims (
    campaign_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    participant_id TEXT NOT NULL,
    claimed_at INTEGER NOT NULL,
    PRIMARY KEY (campaign_id, user_id),
    FOREIGN KEY (campaign_id, participant_id) REFERENCES participants(campaign_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_participant_claims_participant ON participant_claims(participant_id);
