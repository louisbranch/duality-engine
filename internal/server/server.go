// Package server provides the gRPC server for dice rolls.
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/louisbranch/duality-protocol/api/gen/go/duality/v1"
	"github.com/louisbranch/duality-protocol/internal/server/dice"
	"github.com/louisbranch/duality-protocol/internal/server/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server hosts the gRPC dice roll service.
type Server struct {
	pb.UnimplementedDiceRollServiceServer
	listener net.Listener
	grpc     *grpc.Server
	seedFunc func() (int64, error) // Generates per-request random seeds.
}

// New creates a configured gRPC server listening on the provided port.
func New(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("listen on port %d: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{
		listener: listener,
		grpc:     grpcServer,
		seedFunc: random.NewSeed,
	}
	pb.RegisterDiceRollServiceServer(grpcServer, server)

	return server, nil
}

// Serve starts the gRPC server and blocks until it stops.
func (s *Server) Serve() error {
	log.Printf("server listening at %v", s.listener.Addr())
	if err := s.grpc.Serve(s.listener); err != nil {
		return fmt.Errorf("serve gRPC: %w", err)
	}
	return nil
}

// ActionRoll handles action roll requests.
func (s *Server) ActionRoll(ctx context.Context, in *pb.ActionRollRequest) (*pb.ActionRollResponse, error) {
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "action roll request is required")
	}
	if s.seedFunc == nil {
		return nil, status.Error(codes.Internal, "seed generator is not configured")
	}

	// TODO: Expose the seed in the gRPC request/response once the API supports it.
	seed, err := s.seedFunc()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate seed: %v", err)
	}

	var difficulty *int
	if in.Difficulty != nil {
		value := int(*in.Difficulty)
		difficulty = &value
	}

	result, err := dice.RollAction(dice.ActionRequest{
		Modifier:   int(in.GetModifier()),
		Difficulty: difficulty,
		Seed:       seed,
	})
	if err != nil {
		if errors.Is(err, dice.ErrInvalidDifficulty) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to roll action: %v", err)
	}

	response := &pb.ActionRollResponse{
		Duality: &pb.DualityDice{
			HopeD12: int32(result.Hope),
			FearD12: int32(result.Fear),
		},
		Total:   int32(result.Total),
		Outcome: outcomeToProto(result.Outcome),
	}
	if in.Difficulty != nil {
		response.Difficulty = in.Difficulty
	}

	return response, nil
}

// outcomeToProto maps dice outcomes to the protobuf outcome enum.
func outcomeToProto(outcome dice.Outcome) pb.Outcome {
	switch outcome {
	case dice.OutcomeRollWithHope:
		return pb.Outcome_ROLL_WITH_HOPE
	case dice.OutcomeRollWithFear:
		return pb.Outcome_ROLL_WITH_FEAR
	case dice.OutcomeSuccessWithHope:
		return pb.Outcome_SUCCESS_WITH_HOPE
	case dice.OutcomeSuccessWithFear:
		return pb.Outcome_SUCCESS_WITH_FEAR
	case dice.OutcomeFailureWithHope:
		return pb.Outcome_FAILURE_WITH_HOPE
	case dice.OutcomeFailureWithFear:
		return pb.Outcome_FAILURE_WITH_FEAR
	case dice.OutcomeCriticalSuccess:
		return pb.Outcome_CRITICAL_SUCCESS
	default:
		return pb.Outcome_OUTCOME_UNSPECIFIED
	}
}
