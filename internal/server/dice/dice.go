// Package dice implements the dice-rolling logic for the Duality Engine.
package dice

import (
	"errors"
	"math/rand"
)

// Outcome represents the outcome of an action roll.
type Outcome int

const (
	OutcomeUnspecified Outcome = iota
	OutcomeRollWithHope
	OutcomeRollWithFear
	OutcomeSuccessWithHope
	OutcomeSuccessWithFear
	OutcomeFailureWithHope
	OutcomeFailureWithFear
	OutcomeCriticalSuccess
)

func (o Outcome) String() string {
	switch o {
	case OutcomeUnspecified:
		return "Unspecified"
	case OutcomeRollWithHope:
		return "Roll with hope"
	case OutcomeRollWithFear:
		return "Roll with fear"
	case OutcomeSuccessWithHope:
		return "Success with hope"
	case OutcomeSuccessWithFear:
		return "Success with fear"
	case OutcomeFailureWithHope:
		return "Failure with hope"
	case OutcomeFailureWithFear:
		return "Failure with fear"
	case OutcomeCriticalSuccess:
		return "Critical success"
	default:
		return "Unknown"
	}
}

// ErrInvalidDifficulty indicates the difficulty is invalid for a roll.
var ErrInvalidDifficulty = errors.New("difficulty must be non-negative")

// ActionRequest describes an action roll request.
type ActionRequest struct {
	Modifier   int
	Difficulty *int
	Seed       int64
}

// ActionResult contains the outcome of an action roll.
type ActionResult struct {
	Hope    int
	Fear    int
	Total   int
	Outcome Outcome
}

// RollAction performs an action roll from the provided request.
func RollAction(request ActionRequest) (ActionResult, error) {
	if request.Difficulty != nil && *request.Difficulty < 0 {
		return ActionResult{}, ErrInvalidDifficulty
	}

	rng := rand.New(rand.NewSource(request.Seed))
	hope := rollD12(rng)
	fear := rollD12(rng)
	total := hope + fear + request.Modifier

	return ActionResult{
		Hope:    hope,
		Fear:    fear,
		Total:   total,
		Outcome: outcomeFor(hope, fear, total, request.Difficulty),
	}, nil
}

// rollD12 rolls a d12 and returns the result.
func rollD12(rng *rand.Rand) int {
	return rng.Intn(12) + 1
}

// outcomeFor determines the roll outcome based on totals and difficulty.
func outcomeFor(hope int, fear int, total int, difficulty *int) Outcome {
	unopposed := difficulty == nil
	success := !unopposed && total >= *difficulty

	switch {
	case hope == fear:
		return OutcomeCriticalSuccess
	case hope > fear && unopposed:
		return OutcomeRollWithHope
	case fear > hope && unopposed:
		return OutcomeRollWithFear
	case hope > fear && success:
		return OutcomeSuccessWithHope
	case hope < fear && success:
		return OutcomeSuccessWithFear
	case hope > fear:
		return OutcomeFailureWithHope
	case hope < fear:
		return OutcomeFailureWithFear
	default:
		return OutcomeUnspecified
	}
}
