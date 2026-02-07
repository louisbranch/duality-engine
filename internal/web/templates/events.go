package templates

// EventRow represents an event in the timeline (enhanced version).
type EventRow struct {
	CampaignID  string
	Seq         uint64
	Hash        string
	Type        string
	TypeDisplay string
	Timestamp   string
	SessionID   string
	ActorType   string
	ActorName   string
	EntityType  string
	EntityID    string
	EntityName  string
	Description string
	PayloadJSON string
	Expanded    bool
}

// EventFilterOptions holds the current filter state for event lists.
type EventFilterOptions struct {
	SessionID  string
	EventType  string
	ActorType  string
	EntityType string
	StartDate  string
	EndDate    string
}

// EventLogView holds data for rendering the event log page.
type EventLogView struct {
	CampaignID   string
	CampaignName string
	SessionID    string
	SessionName  string
	Events       []EventRow
	Filters      EventFilterOptions
	NextToken    string
	PrevToken    string
	TotalCount   int32
}

// EventTypeOption represents an option in the event type filter dropdown.
type EventTypeOption struct {
	Value   string
	Label   string
	Current bool
}

// GetEventTypeOptions returns the available event type filter options.
func GetEventTypeOptions(current string) []EventTypeOption {
	types := []struct {
		Value string
		Label string
	}{
		{"", "All Types"},
		{"campaign.created", "Campaign Created"},
		{"session.started", "Session Started"},
		{"session.ended", "Session Ended"},
		{"character.created", "Character Created"},
		{"participant.joined", "Participant Joined"},
		{"action.roll_resolved", "Roll Made"},
		{"action.outcome_applied", "Action Taken"},
	}

	options := make([]EventTypeOption, len(types))
	for i, t := range types {
		options[i] = EventTypeOption{
			Value:   t.Value,
			Label:   t.Label,
			Current: t.Value == current,
		}
	}
	return options
}

// GetActorTypeOptions returns the available actor type filter options.
func GetActorTypeOptions(current string) []EventTypeOption {
	types := []struct {
		Value string
		Label string
	}{
		{"", "All Actors"},
		{"system", "System"},
		{"participant", "Participant"},
		{"gm", "GM"},
	}

	options := make([]EventTypeOption, len(types))
	for i, t := range types {
		options[i] = EventTypeOption{
			Value:   t.Value,
			Label:   t.Label,
			Current: t.Value == current,
		}
	}
	return options
}

// GetEntityTypeOptions returns the available entity type filter options.
func GetEntityTypeOptions(current string) []EventTypeOption {
	types := []struct {
		Value string
		Label string
	}{
		{"", "All Entities"},
		{"character", "Character"},
		{"session", "Session"},
		{"campaign", "Campaign"},
		{"participant", "Participant"},
	}

	options := make([]EventTypeOption, len(types))
	for i, t := range types {
		options[i] = EventTypeOption{
			Value:   t.Value,
			Label:   t.Label,
			Current: t.Value == current,
		}
	}
	return options
}
