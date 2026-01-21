package dice

import (
	"testing"
)

func TestNewAction(t *testing.T) {
	action := NewAction(nil)
	if action.rng == nil {
		t.Errorf("NewAction(nil) = %v, want non-nil", action)
	}

	seed := int64(42)
	action = NewAction(&seed)
	if action.rng == nil {
		t.Errorf("NewAction(%d) = %v, want non-nil", seed, action)
	}
}

func TestActionRoll(t *testing.T) {

	diff := func(d int) *int {
		return &d
	}

	tcs := []struct {
		wantOutcome Outcome
		seed        int64
		modifier    int
		difficulty  *int
		wantHope    int
		wantFear    int
		wantTotal   int
	}{
		{
			wantOutcome: OutcomeCriticalSuccess,
			seed:        0,
			modifier:    0,
			difficulty:  nil,
			wantHope:    7,
			wantFear:    7,
			wantTotal:   14,
		},
		{
			wantOutcome: OutcomeRollWithHope,
			seed:        1,
			modifier:    0,
			difficulty:  nil,
			wantHope:    6,
			wantFear:    4,
			wantTotal:   10,
		},
		{
			wantOutcome: OutcomeRollWithFear,
			seed:        3,
			modifier:    0,
			difficulty:  nil,
			wantHope:    5,
			wantFear:    6,
			wantTotal:   11,
		},
		{
			wantOutcome: OutcomeSuccessWithHope,
			seed:        1,
			modifier:    0,
			difficulty:  diff(9),
			wantHope:    6,
			wantFear:    4,
			wantTotal:   10,
		},
		{
			wantOutcome: OutcomeSuccessWithFear,
			seed:        3,
			modifier:    0,
			difficulty:  diff(10),
			wantHope:    5,
			wantFear:    6,
			wantTotal:   11,
		},
		{
			wantOutcome: OutcomeFailureWithHope,
			seed:        1,
			modifier:    -1,
			difficulty:  diff(10),
			wantHope:    6,
			wantFear:    4,
			wantTotal:   9,
		},
		{
			wantOutcome: OutcomeFailureWithFear,
			seed:        3,
			modifier:    -2,
			difficulty:  diff(10),
			wantHope:    5,
			wantFear:    6,
			wantTotal:   9,
		},
	}

	for _, tc := range tcs {
		action := NewAction(&tc.seed)
		gotHope, gotFear, gotTotal, gotOutcome := action.Roll(tc.modifier, tc.difficulty)
		if gotHope != tc.wantHope || gotFear != tc.wantFear || gotTotal != tc.wantTotal || gotOutcome != tc.wantOutcome {
			t.Errorf("Roll(%d, %v) = (%d, %d, %d, %v), want (%d, %d, %d, %v)", tc.modifier, tc.difficulty, gotHope, gotFear, gotTotal, gotOutcome, tc.wantHope, tc.wantFear, tc.wantTotal, tc.wantOutcome)
		}
	}

}
