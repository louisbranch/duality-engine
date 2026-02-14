package game

import (
	"github.com/louisbranch/fracturing.space/internal/services/game/storage"
)

// Stores groups all campaign-related storage interfaces for service injection.
type Stores struct {
	Campaign           storage.CampaignStore
	Participant        storage.ParticipantStore
	ClaimIndex         storage.ClaimIndexStore
	Invite             storage.InviteStore
	Character          storage.CharacterStore
	Daggerheart        storage.DaggerheartStore
	Session            storage.SessionStore
	SessionGate        storage.SessionGateStore
	SessionSpotlight   storage.SessionSpotlightStore
	Event              storage.EventStore
	Telemetry          storage.TelemetryStore
	Statistics         storage.StatisticsStore
	Outcome            storage.RollOutcomeStore
	Snapshot           storage.SnapshotStore
	CampaignFork       storage.CampaignForkStore
	DaggerheartContent storage.DaggerheartContentStore
}
