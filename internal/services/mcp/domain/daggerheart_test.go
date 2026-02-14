package domain

import (
	"context"
	"fmt"
	"testing"

	commonv1 "github.com/louisbranch/fracturing.space/api/gen/go/common/v1"
	pb "github.com/louisbranch/fracturing.space/api/gen/go/systems/daggerheart/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestActionRollHandler(t *testing.T) {
	t.Run("success without rng", func(t *testing.T) {
		difficulty := int32(15)
		client := &fakeDaggerheartClient{
			actionRollResp: &pb.ActionRollResponse{
				Hope:            4,
				Fear:            3,
				Modifier:        2,
				Total:           9,
				IsCrit:          false,
				MeetsDifficulty: false,
				Outcome:         pb.Outcome_FAILURE_WITH_FEAR,
				Difficulty:      &difficulty,
			},
		}
		handler := ActionRollHandler(client)
		_, result, err := handler(context.Background(), nil, ActionRollInput{Modifier: 2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Hope != 4 {
			t.Errorf("expected hope 4, got %d", result.Hope)
		}
		if result.Total != 9 {
			t.Errorf("expected total 9, got %d", result.Total)
		}
		if result.Outcome != "FAILURE_WITH_FEAR" {
			t.Errorf("expected outcome FAILURE_WITH_FEAR, got %q", result.Outcome)
		}
	})

	t.Run("success with rng", func(t *testing.T) {
		seed := uint64(42)
		client := &fakeDaggerheartClient{
			actionRollResp: &pb.ActionRollResponse{
				Hope:     5,
				Fear:     2,
				Modifier: 1,
				Total:    8,
				Outcome:  pb.Outcome_SUCCESS_WITH_HOPE,
				Rng: &commonv1.RngResponse{
					SeedUsed:   42,
					RngAlgo:    "pcg",
					SeedSource: "CLIENT",
					RollMode:   commonv1.RollMode_REPLAY,
				},
			},
		}
		handler := ActionRollHandler(client)
		_, result, err := handler(context.Background(), nil, ActionRollInput{
			Modifier: 1,
			Rng:      &RngRequest{Seed: &seed, RollMode: "REPLAY"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Rng == nil {
			t.Fatal("expected non-nil rng result")
		}
		if result.Rng.SeedUsed != 42 {
			t.Errorf("expected seed 42, got %d", result.Rng.SeedUsed)
		}
		if result.Rng.RollMode != "REPLAY" {
			t.Errorf("expected roll mode REPLAY, got %q", result.Rng.RollMode)
		}
	})

	t.Run("with difficulty", func(t *testing.T) {
		difficulty := 15
		client := &fakeDaggerheartClient{
			actionRollResp: &pb.ActionRollResponse{
				Hope: 6, Fear: 6, Modifier: 3, Total: 15,
				IsCrit: true, MeetsDifficulty: true,
				Outcome: pb.Outcome_CRITICAL_SUCCESS,
			},
		}
		handler := ActionRollHandler(client)
		_, result, err := handler(context.Background(), nil, ActionRollInput{
			Modifier:   3,
			Difficulty: &difficulty,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsCrit {
			t.Error("expected is_crit to be true")
		}
	})
}

func TestDualityOutcomeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			outcomeResp: &pb.DualityOutcomeResponse{
				Hope: 5, Fear: 3, Modifier: 2, Total: 10,
				Outcome:         pb.Outcome_SUCCESS_WITH_HOPE,
				MeetsDifficulty: true,
			},
		}
		handler := DualityOutcomeHandler(client)
		_, result, err := handler(context.Background(), nil, DualityOutcomeInput{
			Hope: 5, Fear: 3, Modifier: 2,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 10 {
			t.Errorf("expected total 10, got %d", result.Total)
		}
		if result.Outcome != "SUCCESS_WITH_HOPE" {
			t.Errorf("expected outcome SUCCESS_WITH_HOPE, got %q", result.Outcome)
		}
	})

	t.Run("with difficulty", func(t *testing.T) {
		difficulty := 12
		respDiff := int32(12)
		client := &fakeDaggerheartClient{
			outcomeResp: &pb.DualityOutcomeResponse{
				Hope: 3, Fear: 4, Modifier: 1, Total: 8,
				Outcome:    pb.Outcome_FAILURE_WITH_FEAR,
				Difficulty: &respDiff,
			},
		}
		handler := DualityOutcomeHandler(client)
		_, result, err := handler(context.Background(), nil, DualityOutcomeInput{
			Hope: 3, Fear: 4, Modifier: 1, Difficulty: &difficulty,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Difficulty == nil || *result.Difficulty != 12 {
			t.Errorf("expected difficulty 12, got %v", result.Difficulty)
		}
	})
}

func TestActionRollHandler_gRPCError(t *testing.T) {
	client := &fakeDaggerheartClient{actionRollErr: fmt.Errorf("error")}
	handler := ActionRollHandler(client)
	_, _, err := handler(context.Background(), nil, ActionRollInput{Modifier: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDualityOutcomeHandler_gRPCError(t *testing.T) {
	client := &fakeDaggerheartClient{outcomeErr: fmt.Errorf("error")}
	handler := DualityOutcomeHandler(client)
	_, _, err := handler(context.Background(), nil, DualityOutcomeInput{Hope: 4, Fear: 3})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDualityExplainHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		stepData, _ := structpb.NewStruct(map[string]any{"key": "value"})
		client := &fakeDaggerheartClient{
			explainResp: &pb.DualityExplainResponse{
				Hope: 4, Fear: 3, Modifier: 2, Total: 9,
				Outcome:      pb.Outcome_SUCCESS_WITH_HOPE,
				RulesVersion: "1.0.0",
				Intermediates: &pb.Intermediates{
					BaseTotal:  7,
					Total:      9,
					HopeGtFear: true,
				},
				Steps: []*pb.ExplainStep{
					{Code: "roll", Message: "Rolled dice", Data: stepData},
				},
			},
		}
		handler := DualityExplainHandler(client)
		_, result, err := handler(context.Background(), nil, DualityExplainInput{
			Hope: 4, Fear: 3, Modifier: 2,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.RulesVersion != "1.0.0" {
			t.Errorf("expected rules_version %q, got %q", "1.0.0", result.RulesVersion)
		}
		if result.Intermediates.BaseTotal != 7 {
			t.Errorf("expected base_total 7, got %d", result.Intermediates.BaseTotal)
		}
		if len(result.Steps) != 1 {
			t.Fatalf("expected 1 step, got %d", len(result.Steps))
		}
		if result.Steps[0].Code != "roll" {
			t.Errorf("expected step code %q, got %q", "roll", result.Steps[0].Code)
		}
	})

	t.Run("nil intermediates", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			explainResp: &pb.DualityExplainResponse{
				Hope: 4, Fear: 3, Total: 9,
				Outcome: pb.Outcome_SUCCESS_WITH_HOPE,
			},
		}
		handler := DualityExplainHandler(client)
		_, _, err := handler(context.Background(), nil, DualityExplainInput{Hope: 4, Fear: 3})
		if err == nil {
			t.Fatal("expected error for nil intermediates")
		}
	})

	t.Run("gRPC error", func(t *testing.T) {
		client := &fakeDaggerheartClient{explainErr: fmt.Errorf("error")}
		handler := DualityExplainHandler(client)
		_, _, err := handler(context.Background(), nil, DualityExplainInput{Hope: 4, Fear: 3})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("nil response", func(t *testing.T) {
		client := &fakeDaggerheartClient{}
		handler := DualityExplainHandler(client)
		_, _, err := handler(context.Background(), nil, DualityExplainInput{Hope: 4, Fear: 3})
		if err == nil {
			t.Fatal("expected error for nil response")
		}
	})

	t.Run("with difficulty", func(t *testing.T) {
		difficulty := 10
		respDiff := int32(10)
		client := &fakeDaggerheartClient{
			explainResp: &pb.DualityExplainResponse{
				Hope: 4, Fear: 3, Modifier: 2, Total: 9,
				Difficulty:   &respDiff,
				Outcome:      pb.Outcome_FAILURE_WITH_FEAR,
				RulesVersion: "1.0.0",
				Intermediates: &pb.Intermediates{
					BaseTotal:  7,
					Total:      9,
					HopeGtFear: true,
				},
			},
		}
		handler := DualityExplainHandler(client)
		_, result, err := handler(context.Background(), nil, DualityExplainInput{
			Hope: 4, Fear: 3, Modifier: 2, Difficulty: &difficulty,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Difficulty == nil || *result.Difficulty != 10 {
			t.Errorf("expected difficulty 10, got %v", result.Difficulty)
		}
	})

	t.Run("with request_id", func(t *testing.T) {
		reqID := "req-123"
		client := &fakeDaggerheartClient{
			explainResp: &pb.DualityExplainResponse{
				Hope: 4, Fear: 3, Total: 9,
				Outcome:      pb.Outcome_SUCCESS_WITH_HOPE,
				RulesVersion: "1.0.0",
				Intermediates: &pb.Intermediates{
					BaseTotal:  7,
					Total:      9,
					HopeGtFear: true,
				},
			},
		}
		handler := DualityExplainHandler(client)
		_, _, err := handler(context.Background(), nil, DualityExplainInput{
			Hope: 4, Fear: 3, RequestID: &reqID,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("step with nil data", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			explainResp: &pb.DualityExplainResponse{
				Hope: 4, Fear: 3, Total: 9,
				Outcome:      pb.Outcome_SUCCESS_WITH_HOPE,
				RulesVersion: "1.0.0",
				Intermediates: &pb.Intermediates{
					BaseTotal: 7, Total: 9,
				},
				Steps: []*pb.ExplainStep{
					{Code: "roll", Message: "Rolled dice"},
				},
			},
		}
		handler := DualityExplainHandler(client)
		_, result, err := handler(context.Background(), nil, DualityExplainInput{Hope: 4, Fear: 3})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Steps) != 1 {
			t.Fatalf("expected 1 step, got %d", len(result.Steps))
		}
		if result.Steps[0].Data == nil {
			t.Error("expected non-nil data map")
		}
	})
}

func TestDualityProbabilityHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			probabilityResp: &pb.DualityProbabilityResponse{
				TotalOutcomes: 36,
				CritCount:     6,
				SuccessCount:  15,
				FailureCount:  15,
				OutcomeCounts: []*pb.OutcomeCount{
					{Outcome: pb.Outcome_CRITICAL_SUCCESS, Count: 6},
					{Outcome: pb.Outcome_SUCCESS_WITH_HOPE, Count: 10},
					{Outcome: pb.Outcome_SUCCESS_WITH_FEAR, Count: 5},
					{Outcome: pb.Outcome_FAILURE_WITH_HOPE, Count: 5},
					{Outcome: pb.Outcome_FAILURE_WITH_FEAR, Count: 10},
				},
			},
		}
		handler := DualityProbabilityHandler(client)
		_, result, err := handler(context.Background(), nil, DualityProbabilityInput{
			Modifier: 2, Difficulty: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.TotalOutcomes != 36 {
			t.Errorf("expected 36 total outcomes, got %d", result.TotalOutcomes)
		}
		if len(result.OutcomeCounts) != 5 {
			t.Errorf("expected 5 outcome counts, got %d", len(result.OutcomeCounts))
		}
	})
}

func TestDualityProbabilityHandler_gRPCError(t *testing.T) {
	client := &fakeDaggerheartClient{probabilityErr: fmt.Errorf("error")}
	handler := DualityProbabilityHandler(client)
	_, _, err := handler(context.Background(), nil, DualityProbabilityInput{Modifier: 2, Difficulty: 10})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDualityProbabilityHandler_nilResponse(t *testing.T) {
	client := &fakeDaggerheartClient{}
	handler := DualityProbabilityHandler(client)
	_, _, err := handler(context.Background(), nil, DualityProbabilityInput{Modifier: 2, Difficulty: 10})
	if err == nil {
		t.Fatal("expected error for nil response")
	}
}

func TestRulesVersionHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			rulesVersionResp: &pb.RulesVersionResponse{
				System:         "daggerheart",
				Module:         "duality",
				RulesVersion:   "1.0.0",
				DiceModel:      "2d12",
				TotalFormula:   "hope + fear + modifier",
				CritRule:       "hope == fear",
				DifficultyRule: "total >= difficulty",
				Outcomes: []pb.Outcome{
					pb.Outcome_CRITICAL_SUCCESS,
					pb.Outcome_SUCCESS_WITH_HOPE,
				},
			},
		}
		handler := RulesVersionHandler(client)
		_, result, err := handler(context.Background(), nil, RulesVersionInput{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.System != "daggerheart" {
			t.Errorf("expected system %q, got %q", "daggerheart", result.System)
		}
		if result.RulesVersion != "1.0.0" {
			t.Errorf("expected rules_version %q, got %q", "1.0.0", result.RulesVersion)
		}
		if len(result.Outcomes) != 2 {
			t.Errorf("expected 2 outcomes, got %d", len(result.Outcomes))
		}
	})
}

func TestRulesVersionHandler_gRPCError(t *testing.T) {
	client := &fakeDaggerheartClient{rulesVersionErr: fmt.Errorf("error")}
	handler := RulesVersionHandler(client)
	_, _, err := handler(context.Background(), nil, RulesVersionInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRulesVersionHandler_nilResponse(t *testing.T) {
	client := &fakeDaggerheartClient{}
	handler := RulesVersionHandler(client)
	_, _, err := handler(context.Background(), nil, RulesVersionInput{})
	if err == nil {
		t.Fatal("expected error for nil response")
	}
}

func TestRollDiceHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := &fakeDaggerheartClient{
			rollDiceResp: &pb.RollDiceResponse{
				Rolls: []*pb.DiceRoll{
					{Sides: 6, Results: []int32{3, 4}, Total: 7},
					{Sides: 12, Results: []int32{8}, Total: 8},
				},
				Total: 15,
			},
		}
		handler := RollDiceHandler(client)
		_, result, err := handler(context.Background(), nil, RollDiceInput{
			Dice: []RollDiceSpec{
				{Sides: 6, Count: 2},
				{Sides: 12, Count: 1},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Total != 15 {
			t.Errorf("expected total 15, got %d", result.Total)
		}
		if len(result.Rolls) != 2 {
			t.Fatalf("expected 2 rolls, got %d", len(result.Rolls))
		}
		if result.Rolls[0].Sides != 6 {
			t.Errorf("expected sides 6, got %d", result.Rolls[0].Sides)
		}
	})

	t.Run("with rng", func(t *testing.T) {
		seed := uint64(99)
		client := &fakeDaggerheartClient{
			rollDiceResp: &pb.RollDiceResponse{
				Rolls: []*pb.DiceRoll{{Sides: 6, Results: []int32{1}, Total: 1}},
				Total: 1,
				Rng: &commonv1.RngResponse{
					SeedUsed:   99,
					RngAlgo:    "pcg",
					SeedSource: "CLIENT",
					RollMode:   commonv1.RollMode_REPLAY,
				},
			},
		}
		handler := RollDiceHandler(client)
		_, result, err := handler(context.Background(), nil, RollDiceInput{
			Dice: []RollDiceSpec{{Sides: 6, Count: 1}},
			Rng:  &RngRequest{Seed: &seed, RollMode: "REPLAY"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Rng == nil {
			t.Fatal("expected non-nil rng")
		}
		if result.Rng.SeedUsed != 99 {
			t.Errorf("expected seed 99, got %d", result.Rng.SeedUsed)
		}
	})

	t.Run("gRPC error", func(t *testing.T) {
		client := &fakeDaggerheartClient{rollDiceErr: fmt.Errorf("error")}
		handler := RollDiceHandler(client)
		_, _, err := handler(context.Background(), nil, RollDiceInput{Dice: []RollDiceSpec{{Sides: 6, Count: 1}}})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("nil response", func(t *testing.T) {
		client := &fakeDaggerheartClient{}
		handler := RollDiceHandler(client)
		_, _, err := handler(context.Background(), nil, RollDiceInput{Dice: []RollDiceSpec{{Sides: 6, Count: 1}}})
		if err == nil {
			t.Fatal("expected error for nil response")
		}
	})
}
