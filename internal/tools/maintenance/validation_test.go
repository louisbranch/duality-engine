package maintenance

import (
	"encoding/json"
	"testing"

	"github.com/louisbranch/fracturing.space/internal/services/game/domain/campaign/event"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/core/dice"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/systems/daggerheart"
)

func TestIsSnapshotEvent(t *testing.T) {
	tests := []struct {
		name     string
		systemID string
		want     bool
	}{
		{"empty system id", "", false},
		{"whitespace only", "   ", false},
		{"daggerheart", "GAME_SYSTEM_DAGGERHEART", true},
		{"any system", "custom", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			evt := event.Event{SystemID: tc.systemID}
			if got := isSnapshotEvent(evt); got != tc.want {
				t.Errorf("isSnapshotEvent(systemID=%q) = %v, want %v", tc.systemID, got, tc.want)
			}
		})
	}
}

func intPtr(v int) *int { return &v }

func TestValidateSnapshotEvent_CharacterStatePatched(t *testing.T) {
	makeEvent := func(payload daggerheart.CharacterStatePatchedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeCharacterStatePatched,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{
			CharacterID: "char-1",
			HpAfter:     intPtr(5),
			HopeAfter:   intPtr(2),
			StressAfter: intPtr(3),
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing character id", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for missing character id")
		}
	})

	t.Run("hp out of range", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{
			CharacterID: "char-1",
			HpAfter:     intPtr(-1),
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for negative HP")
		}
	})

	t.Run("hope out of range", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{
			CharacterID:  "char-1",
			HopeMaxAfter: intPtr(100),
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for hope_max out of range")
		}
	})

	t.Run("stress out of range", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{
			CharacterID: "char-1",
			StressAfter: intPtr(-1),
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for negative stress")
		}
	})

	t.Run("armor out of range", func(t *testing.T) {
		evt := makeEvent(daggerheart.CharacterStatePatchedPayload{
			CharacterID: "char-1",
			ArmorAfter:  intPtr(-1),
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for negative armor")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		evt := event.Event{
			Type:        daggerheart.EventTypeCharacterStatePatched,
			PayloadJSON: []byte(`{invalid`),
		}
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestValidateSnapshotEvent_DeathMoveResolved(t *testing.T) {
	makeEvent := func(payload daggerheart.DeathMoveResolvedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeDeathMoveResolved,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.DeathMoveResolvedPayload{
			CharacterID:    "char-1",
			Move:           daggerheart.DeathMoveAvoidDeath,
			LifeStateAfter: daggerheart.LifeStateDead,
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing character id", func(t *testing.T) {
		evt := makeEvent(daggerheart.DeathMoveResolvedPayload{
			Move:           daggerheart.DeathMoveAvoidDeath,
			LifeStateAfter: daggerheart.LifeStateDead,
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for missing character id")
		}
	})

	t.Run("missing life state", func(t *testing.T) {
		evt := makeEvent(daggerheart.DeathMoveResolvedPayload{
			CharacterID: "char-1",
			Move:        daggerheart.DeathMoveAvoidDeath,
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for missing life_state_after")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		evt := event.Event{
			Type:        daggerheart.EventTypeDeathMoveResolved,
			PayloadJSON: []byte(`{bad`),
		}
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestValidateSnapshotEvent_BlazeOfGloryResolved(t *testing.T) {
	makeEvent := func(payload daggerheart.BlazeOfGloryResolvedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeBlazeOfGloryResolved,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.BlazeOfGloryResolvedPayload{
			CharacterID:    "char-1",
			LifeStateAfter: daggerheart.LifeStateDead,
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing character id", func(t *testing.T) {
		evt := makeEvent(daggerheart.BlazeOfGloryResolvedPayload{
			LifeStateAfter: daggerheart.LifeStateDead,
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("missing life state", func(t *testing.T) {
		evt := makeEvent(daggerheart.BlazeOfGloryResolvedPayload{
			CharacterID: "char-1",
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})
}

func TestValidateSnapshotEvent_AttackResolved(t *testing.T) {
	makeEvent := func(payload daggerheart.AttackResolvedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeAttackResolved,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			CharacterID: "char-1",
			RollSeq:     1,
			Targets:     []string{"goblin-1"},
			Outcome:     "ROLL_WITH_HOPE",
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing character id", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			RollSeq: 1, Targets: []string{"t"}, Outcome: "o",
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("missing roll seq", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			CharacterID: "c", Targets: []string{"t"}, Outcome: "o",
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("empty targets", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			CharacterID: "c", RollSeq: 1, Outcome: "o",
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("empty target value", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			CharacterID: "c", RollSeq: 1, Targets: []string{""}, Outcome: "o",
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("missing outcome", func(t *testing.T) {
		evt := makeEvent(daggerheart.AttackResolvedPayload{
			CharacterID: "c", RollSeq: 1, Targets: []string{"t"},
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})
}

func TestValidateSnapshotEvent_ReactionResolved(t *testing.T) {
	makeEvent := func(payload daggerheart.ReactionResolvedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeReactionResolved,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.ReactionResolvedPayload{
			CharacterID: "char-1", RollSeq: 1, Outcome: "ROLL_WITH_FEAR",
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		evt := makeEvent(daggerheart.ReactionResolvedPayload{})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})
}

func TestValidateSnapshotEvent_DamageRollResolved(t *testing.T) {
	makeEvent := func(payload daggerheart.DamageRollResolvedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeDamageRollResolved,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.DamageRollResolvedPayload{
			CharacterID: "char-1",
			RollSeq:     1,
			Rolls:       []dice.Roll{{Sides: 6, Results: []int{3}, Total: 3}},
		})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing rolls", func(t *testing.T) {
		evt := makeEvent(daggerheart.DamageRollResolvedPayload{
			CharacterID: "char-1", RollSeq: 1,
		})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error")
		}
	})
}

func TestValidateSnapshotEvent_GMFearChanged(t *testing.T) {
	makeEvent := func(payload daggerheart.GMFearChangedPayload) event.Event {
		data, _ := json.Marshal(payload)
		return event.Event{
			Type:        daggerheart.EventTypeGMFearChanged,
			PayloadJSON: data,
		}
	}

	t.Run("valid", func(t *testing.T) {
		evt := makeEvent(daggerheart.GMFearChangedPayload{After: 3})
		if err := validateSnapshotEvent(evt); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("out of range", func(t *testing.T) {
		evt := makeEvent(daggerheart.GMFearChangedPayload{After: -1})
		if err := validateSnapshotEvent(evt); err == nil {
			t.Error("expected error for negative fear")
		}
	})
}

func TestValidateSnapshotEvent_UnknownType(t *testing.T) {
	evt := event.Event{Type: "unknown.event", PayloadJSON: []byte("{}")}
	if err := validateSnapshotEvent(evt); err != nil {
		t.Errorf("unknown event type should not error: %v", err)
	}
}
