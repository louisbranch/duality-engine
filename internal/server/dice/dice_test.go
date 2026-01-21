package dice

import (
	"errors"
	"testing"
)

func TestRollAction(t *testing.T) {
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
		result, err := RollAction(ActionRequest{
			Modifier:   tc.modifier,
			Difficulty: tc.difficulty,
			Seed:       tc.seed,
		})
		if err != nil {
			t.Fatalf("RollAction returned error: %v", err)
		}
		if result.Hope != tc.wantHope || result.Fear != tc.wantFear || result.Total != tc.wantTotal || result.Outcome != tc.wantOutcome {
			t.Errorf("RollAction(%d, %v) = (%d, %d, %d, %v), want (%d, %d, %d, %v)", tc.modifier, tc.difficulty, result.Hope, result.Fear, result.Total, result.Outcome, tc.wantHope, tc.wantFear, tc.wantTotal, tc.wantOutcome)
		}
	}
}

func TestRollActionRejectsNegativeDifficulty(t *testing.T) {
	difficulty := -1
	_, err := RollAction(ActionRequest{
		Modifier:   0,
		Difficulty: &difficulty,
		Seed:       0,
	})
	if !errors.Is(err, ErrInvalidDifficulty) {
		t.Fatalf("RollAction error = %v, want %v", err, ErrInvalidDifficulty)
	}
}
