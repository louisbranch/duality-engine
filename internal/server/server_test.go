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
