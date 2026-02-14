package scenario

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	gamev1 "github.com/louisbranch/fracturing.space/api/gen/go/game/v1"
	daggerheartv1 "github.com/louisbranch/fracturing.space/api/gen/go/systems/daggerheart/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config controls scenario execution.
type Config struct {
	GRPCAddr   string
	Timeout    time.Duration
	Assertions AssertionMode
	Verbose    bool
	Logger     *log.Logger
}

// DefaultConfig returns default runner configuration.
func DefaultConfig() Config {
	return Config{
		GRPCAddr:   "localhost:8080",
		Timeout:    10 * time.Second,
		Assertions: AssertionStrict,
		Verbose:    false,
	}
}

// Runner executes Lua scenarios against the game gRPC API.
type Runner struct {
	conn       *grpc.ClientConn
	env        scenarioEnv
	assertions Assertions
	logger     *log.Logger
	verbose    bool
	timeout    time.Duration
	auth       *MockAuth
	userID     string
}

// NewRunner connects to gRPC and prepares a scenario runner.
func NewRunner(ctx context.Context, cfg Config) (*Runner, error) {
	if cfg.GRPCAddr == "" {
		return nil, errors.New("grpc address is required")
	}
	logger := cfg.Logger
	if logger == nil {
		logger = log.New(os.Stderr, "", 0)
	}
	assertions := Assertions{Mode: cfg.Assertions, Logger: logger}

	conn, err := grpc.NewClient(
		cfg.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	)
	if err != nil {
		return nil, fmt.Errorf("dial gRPC: %w", err)
	}

	auth := NewMockAuth()
	userID := auth.CreateUser("Scenario Runner")
	if userID == "" {
		_ = conn.Close()
		return nil, errors.New("mock auth returned empty user id")
	}

	env := scenarioEnv{
		campaignClient:    gamev1.NewCampaignServiceClient(conn),
		participantClient: gamev1.NewParticipantServiceClient(conn),
		sessionClient:     gamev1.NewSessionServiceClient(conn),
		characterClient:   gamev1.NewCharacterServiceClient(conn),
		snapshotClient:    gamev1.NewSnapshotServiceClient(conn),
		eventClient:       gamev1.NewEventServiceClient(conn),
		daggerheartClient: daggerheartv1.NewDaggerheartServiceClient(conn),
		userID:            userID,
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Runner{
		conn:       conn,
		env:        env,
		assertions: assertions,
		logger:     logger,
		verbose:    cfg.Verbose,
		timeout:    timeout,
		auth:       auth,
		userID:     userID,
	}, nil
}

// Close releases resources held by the runner.
func (r *Runner) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// RunFile loads and executes a scenario file.
func RunFile(ctx context.Context, cfg Config, path string) error {
	runner, err := NewRunner(ctx, cfg)
	if err != nil {
		return err
	}
	defer runner.Close()

	scenario, err := LoadScenarioFromFile(path)
	if err != nil {
		return err
	}
	return runner.RunScenario(ctx, scenario)
}

// RunScenario executes the scenario steps against gRPC.
func (r *Runner) RunScenario(ctx context.Context, scenario *Scenario) error {
	if scenario == nil {
		return errors.New("scenario is required")
	}
	r.logf("scenario start: %s (%d steps)", scenario.Name, len(scenario.Steps))
	state := &scenarioState{
		actors:       map[string]string{},
		adversaries:  map[string]string{},
		countdowns:   map[string]string{},
		participants: map[string]string{},
		userID:       r.userID,
	}

	for index, step := range scenario.Steps {
		step := step
		stepNumber := index + 1
		r.logf("step %d/%d start: %s", stepNumber, len(scenario.Steps), step.Kind)
		stepStart := time.Now()
		stepCtx, cancel := context.WithTimeout(ctx, r.timeout)
		err := r.runStep(stepCtx, state, step)
		cancel()
		if err != nil {
			return fmt.Errorf("step %d (%s): %w", index+1, step.Kind, err)
		}
		r.logf("step %d/%d done: %s (%s)", stepNumber, len(scenario.Steps), step.Kind, time.Since(stepStart))
	}
	r.logf("scenario done: %s", scenario.Name)
	return nil
}

func (r *Runner) logf(format string, args ...any) {
	if !r.verbose || r.logger == nil {
		return
	}
	r.logger.Printf(format, args...)
}
