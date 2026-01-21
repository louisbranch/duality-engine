// Package dice implements the dice-rolling logic for the Duality Engine.
package dice

import (
	crand "crypto/rand"
	"encoding/binary"
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

// Action represents a single action roll.
type Action struct {
	rng *rand.Rand
}

// NewAction creates a new Action with an optional random seed.
func NewAction(seed *int64) Action {
	var s int64
	if seed != nil {
		s = *seed
	} else {
		s = defaultSeed()
	}
	return Action{
		rng: rand.New(rand.NewSource(s)),
	}
}

// Roll performs an action roll with the given modifier and difficulty.
// Returns the hope, fear, total, and outcome of the roll.
func (a *Action) Roll(modifier int, difficulty *int) (hope int, fear int, total int, outcome Outcome) {
	hope = a.rollD12()
	fear = a.rollD12()
	total = hope + fear + modifier

	unopposed := difficulty == nil
	sucess := !unopposed && total >= *difficulty

	switch {
	case hope == fear:
		outcome = OutcomeCriticalSuccess
	case hope > fear && unopposed:
		outcome = OutcomeRollWithHope
	case fear > hope && unopposed:
		outcome = OutcomeRollWithFear
	case hope > fear && sucess:
		outcome = OutcomeSuccessWithHope
	case hope < fear && sucess:
		outcome = OutcomeSuccessWithFear
	case hope > fear:
		outcome = OutcomeFailureWithHope
	case hope < fear:
		outcome = OutcomeFailureWithFear
	default:
		// Branch should be unreacheable
		outcome = OutcomeUnspecified
	}

	return
}

// rollD12 rolls a d12 and returns the result.
func (a *Action) rollD12() int {
	return a.rng.Intn(12) + 1
}

// defaultSeed generates a random seed using crypto/rand.
func defaultSeed() int64 {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}
