package daggerheart

import "github.com/louisbranch/fracturing.space/internal/services/game/storage"

// Stores groups storage interfaces used by the Daggerheart service.
type Stores struct {
	Campaign           storage.CampaignStore
	Character          storage.CharacterStore
	Session            storage.SessionStore
	SessionGate        storage.SessionGateStore
	SessionSpotlight   storage.SessionSpotlightStore
	Daggerheart        storage.DaggerheartStore
	DaggerheartContent storage.DaggerheartContentStore
	Event              storage.EventStore
}
