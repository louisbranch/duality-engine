package server

import (
	"context"
	"errors"
	"testing"

	pb "github.com/louisbranch/duality-protocol/api/gen/go/duality/v1"
	"github.com/louisbranch/duality-protocol/internal/server/dice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestActionRollRejectsNilRequest(t *testing.T) {
	server := newTestServer(42)

	_, err := server.ActionRoll(context.Background(), nil)
	assertStatusCode(t, err, codes.InvalidArgument)
}

func TestActionRollRejectsNegativeDifficulty(t *testing.T) {
	server := newTestServer(42)

	negative := int32(-1)
	_, err := server.ActionRoll(context.Background(), &pb.ActionRollRequest{Difficulty: &negative})
	assertStatusCode(t, err, codes.InvalidArgument)
}

func TestActionRollWithDifficulty(t *testing.T) {
	seed := int64(99)
	server := newTestServer(seed)

	difficulty := int32(10)
	modifier := int32(2)
	response, err := server.ActionRoll(context.Background(), &pb.ActionRollRequest{
		Modifier:   modifier,
		Difficulty: &difficulty,
	})
	if err != nil {
		t.Fatalf("ActionRoll returned error: %v", err)
	}
	assertResponseMatches(t, response, seed, modifier, &difficulty)
}

func TestActionRollWithoutDifficulty(t *testing.T) {
	seed := int64(55)
	server := newTestServer(seed)

	modifier := int32(-1)
	response, err := server.ActionRoll(context.Background(), &pb.ActionRollRequest{Modifier: modifier})
	if err != nil {
		t.Fatalf("ActionRoll returned error: %v", err)
	}
	if response.Difficulty != nil {
		t.Fatalf("ActionRoll difficulty = %v, want nil", *response.Difficulty)
	}
	assertResponseMatches(t, response, seed, modifier, nil)
}

func TestActionRollSeedFailure(t *testing.T) {
	server := &Server{
		seedFunc: func() (int64, error) {
			return 0, errors.New("seed failure")
		},
	}

	_, err := server.ActionRoll(context.Background(), &pb.ActionRollRequest{})
	assertStatusCode(t, err, codes.Internal)
}

func TestRollDiceRejectsNilRequest(t *testing.T) {
	server := newTestServer(42)

	_, err := server.RollDice(context.Background(), nil)
	assertStatusCode(t, err, codes.InvalidArgument)
}

func TestRollDiceRejectsMissingDice(t *testing.T) {
	server := newTestServer(42)

	_, err := server.RollDice(context.Background(), &pb.RollDiceRequest{})
	assertStatusCode(t, err, codes.InvalidArgument)
}

func TestRollDiceRejectsInvalidDiceSpec(t *testing.T) {
	server := newTestServer(42)

	_, err := server.RollDice(context.Background(), &pb.RollDiceRequest{
		Dice: []*pb.DiceSpec{{Sides: 0, Count: 1}},
	})
	assertStatusCode(t, err, codes.InvalidArgument)
}

func TestRollDiceReturnsResults(t *testing.T) {
	seed := int64(13)
	server := newTestServer(seed)

	response, err := server.RollDice(context.Background(), &pb.RollDiceRequest{
		Dice: []*pb.DiceSpec{
			{Sides: 6, Count: 2},
			{Sides: 8, Count: 1},
		},
	})
	if err != nil {
		t.Fatalf("RollDice returned error: %v", err)
	}
	assertRollDiceResponse(t, response, seed, []dice.DiceSpec{{Sides: 6, Count: 2}, {Sides: 8, Count: 1}})
}

func TestRollDiceSeedFailure(t *testing.T) {
	server := &Server{
		seedFunc: func() (int64, error) {
			return 0, errors.New("seed failure")
		},
	}

	_, err := server.RollDice(context.Background(), &pb.RollDiceRequest{
		Dice: []*pb.DiceSpec{{Sides: 6, Count: 1}},
	})
	assertStatusCode(t, err, codes.Internal)
}

// assertStatusCode verifies the gRPC status code for an error.
func assertStatusCode(t *testing.T, err error, want codes.Code) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error with code %v", want)
	}
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T", err)
	}
	if statusErr.Code() != want {
		t.Fatalf("status code = %v, want %v", statusErr.Code(), want)
	}
}

// assertResponseMatches validates response fields against expectations.
func assertResponseMatches(t *testing.T, response *pb.ActionRollResponse, seed int64, modifier int32, difficulty *int32) {
	t.Helper()

	if response == nil {
		t.Fatal("ActionRoll response is nil")
	}
	if response.Duality == nil {
		t.Fatal("ActionRoll duality dice is nil")
	}

	result, err := dice.RollAction(dice.ActionRequest{
		Modifier:   int(modifier),
		Difficulty: intPointer(difficulty),
		Seed:       seed,
	})
	if err != nil {
		t.Fatalf("RollAction returned error: %v", err)
	}

	if response.Duality.GetHopeD12() != int32(result.Hope) || response.Duality.GetFearD12() != int32(result.Fear) {
		t.Fatalf("ActionRoll duality dice = (%d, %d), want (%d, %d)", response.Duality.GetHopeD12(), response.Duality.GetFearD12(), result.Hope, result.Fear)
	}
	if response.Total != int32(result.Total) {
		t.Fatalf("ActionRoll total = %d, want %d", response.Total, result.Total)
	}
	if response.Outcome != outcomeToProto(result.Outcome) {
		t.Fatalf("ActionRoll outcome = %v, want %v", response.Outcome, outcomeToProto(result.Outcome))
	}
	if difficulty != nil && response.Difficulty == nil {
		t.Fatal("ActionRoll difficulty is nil, want value")
	}
	if difficulty != nil && response.Difficulty != nil && *response.Difficulty != *difficulty {
		t.Fatalf("ActionRoll difficulty = %d, want %d", *response.Difficulty, *difficulty)
	}
}

// assertRollDiceResponse validates roll dice response fields against expectations.
func assertRollDiceResponse(t *testing.T, response *pb.RollDiceResponse, seed int64, specs []dice.DiceSpec) {
	t.Helper()

	if response == nil {
		t.Fatal("RollDice response is nil")
	}

	result, err := dice.RollDice(dice.RollRequest{
		Dice: specs,
		Seed: seed,
	})
	if err != nil {
		t.Fatalf("RollDice returned error: %v", err)
	}

	if len(response.GetRolls()) != len(result.Rolls) {
		t.Fatalf("RollDice roll count = %d, want %d", len(response.GetRolls()), len(result.Rolls))
	}
	if response.Total != int32(result.Total) {
		t.Fatalf("RollDice total = %d, want %d", response.Total, result.Total)
	}

	for i, roll := range response.GetRolls() {
		want := result.Rolls[i]
		if roll.GetSides() != int32(want.Sides) {
			t.Fatalf("RollDice roll[%d] sides = %d, want %d", i, roll.GetSides(), want.Sides)
		}
		if roll.GetTotal() != int32(want.Total) {
			t.Fatalf("RollDice roll[%d] total = %d, want %d", i, roll.GetTotal(), want.Total)
		}
		if len(roll.GetResults()) != len(want.Results) {
			t.Fatalf("RollDice roll[%d] results = %v, want %v", i, roll.GetResults(), want.Results)
		}
		for j, value := range roll.GetResults() {
			if value != int32(want.Results[j]) {
				t.Fatalf("RollDice roll[%d] result[%d] = %d, want %d", i, j, value, want.Results[j])
			}
		}
	}
}

// newTestServer creates a server with a fixed seed generator.
func newTestServer(seed int64) *Server {
	return &Server{
		seedFunc: func() (int64, error) {
			return seed, nil
		},
	}
}

// intPointer converts a difficulty pointer to the dice package type.
func intPointer(value *int32) *int {
	if value == nil {
		return nil
	}

	converted := int(*value)
	return &converted
}
