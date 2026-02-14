package game

import (
	"context"
	"testing"

	"github.com/louisbranch/fracturing.space/internal/services/game/domain/campaign/event"
	"github.com/louisbranch/fracturing.space/internal/services/game/domain/systems/daggerheart"
)

func TestApplyStressVulnerableCondition_AddsCondition(t *testing.T) {
	ctx := context.Background()
	eventStore := newFakeEventStore()
	dhStore := newFakeDaggerheartStore()

	err := applyStressVulnerableCondition(
		ctx,
		Stores{Event: eventStore, Daggerheart: dhStore},
		"c1",
		"s1",
		"ch1",
		nil,
		2,
		6,
		6,
		event.ActorTypeGM,
		"gm-1",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := len(eventStore.events["c1"]); got != 1 {
		t.Fatalf("expected 1 event, got %d", got)
	}
	if eventStore.events["c1"][0].Type != daggerheart.EventTypeConditionChanged {
		t.Fatalf("event type = %s, want %s", eventStore.events["c1"][0].Type, daggerheart.EventTypeConditionChanged)
	}
	state, err := dhStore.GetDaggerheartCharacterState(ctx, "c1", "ch1")
	if err != nil {
		t.Fatalf("expected daggerheart state, got %v", err)
	}
	if !containsCondition(state.Conditions, daggerheart.ConditionVulnerable) {
		t.Fatalf("expected vulnerable condition, got %v", state.Conditions)
	}
}

func TestApplyStressVulnerableCondition_RemovesCondition(t *testing.T) {
	ctx := context.Background()
	eventStore := newFakeEventStore()
	dhStore := newFakeDaggerheartStore()

	err := applyStressVulnerableCondition(
		ctx,
		Stores{Event: eventStore, Daggerheart: dhStore},
		"c1",
		"s1",
		"ch1",
		[]string{daggerheart.ConditionVulnerable},
		6,
		5,
		6,
		event.ActorTypeGM,
		"gm-1",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := len(eventStore.events["c1"]); got != 1 {
		t.Fatalf("expected 1 event, got %d", got)
	}
	state, err := dhStore.GetDaggerheartCharacterState(ctx, "c1", "ch1")
	if err != nil {
		t.Fatalf("expected daggerheart state, got %v", err)
	}
	if containsCondition(state.Conditions, daggerheart.ConditionVulnerable) {
		t.Fatalf("expected vulnerable condition removed, got %v", state.Conditions)
	}
}

func TestApplyStressVulnerableCondition_NoOpWhenUnchanged(t *testing.T) {
	ctx := context.Background()
	eventStore := newFakeEventStore()
	dhStore := newFakeDaggerheartStore()

	err := applyStressVulnerableCondition(
		ctx,
		Stores{Event: eventStore, Daggerheart: dhStore},
		"c1",
		"s1",
		"ch1",
		nil,
		3,
		3,
		6,
		event.ActorTypeGM,
		"gm-1",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := len(eventStore.events["c1"]); got != 0 {
		t.Fatalf("expected 0 events, got %d", got)
	}
}

func TestApplyStressVulnerableCondition_NoOpWhenAlreadyVulnerable(t *testing.T) {
	ctx := context.Background()
	eventStore := newFakeEventStore()
	dhStore := newFakeDaggerheartStore()

	err := applyStressVulnerableCondition(
		ctx,
		Stores{Event: eventStore, Daggerheart: dhStore},
		"c1",
		"s1",
		"ch1",
		[]string{daggerheart.ConditionVulnerable},
		5,
		6,
		6,
		event.ActorTypeGM,
		"gm-1",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got := len(eventStore.events["c1"]); got != 0 {
		t.Fatalf("expected 0 events, got %d", got)
	}
}

func containsCondition(conditions []string, target string) bool {
	for _, condition := range conditions {
		if condition == target {
			return true
		}
	}
	return false
}
