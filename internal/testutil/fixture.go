package testutil

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/mock"
)

// Fixture provides affordances for testing CLI commands.
type Fixture struct {
	MockController *gomock.Controller
	Client         *mock.Client
	ActionWaiter   *state.MockActionWaiter
	TokenEnsurer   *state.MockTokenEnsurer
	Config         *config.MockConfig
}

// NewFixture creates a new Fixture.
func NewFixture(t *testing.T) *Fixture {
	ctrl := gomock.NewController(t)

	return &Fixture{
		MockController: ctrl,
		Client:         mock.NewClient(ctrl),
		ActionWaiter:   state.NewMockActionWaiter(ctrl),
		TokenEnsurer:   state.NewMockTokenEnsurer(ctrl),
		Config:         config.NewMockConfig(ctrl),
	}
}

// ExpectEnsureToken makes the mock TokenEnsurer expect a EnsureToken call.
func (f *Fixture) ExpectEnsureToken() {
	f.TokenEnsurer.EXPECT().EnsureToken(gomock.Any(), gomock.Any()).Return(nil)
}

// Run runs the given cobra command with the given arguments and returns stdout output and a
// potential error.
func (f *Fixture) Run(cmd *cobra.Command, args []string) (string, string, error) {
	cmd.SetArgs(args)
	return CaptureOutStreams(func() error {
		// We need to re-set the output because CaptureOutStream changes os.Stdout
		// and the command's outWriter still points to the old os.Stdout.
		cmd.SetOut(os.Stdout)
		return cmd.Execute()
	})
}

// Finish must be called after the test is finished, preferably via `defer` directly after
// creating the Fixture.
func (f *Fixture) Finish() {
	f.MockController.Finish()
}

// fixtureState implements state.State for testing purposes.
type fixtureState struct {
	context.Context
	state.TokenEnsurer
	state.ActionWaiter

	client hcapi2.Client
	config config.Config
}

func (*fixtureState) WriteConfig() error {
	return errors.New("not implemented")
}

func (s *fixtureState) Client() hcapi2.Client {
	return s.client
}

func (s *fixtureState) Config() config.Config {
	return s.config
}

// State returns a state.State implementation for testing purposes.
func (f *Fixture) State() state.State {
	return &fixtureState{
		Context:      context.Background(),
		TokenEnsurer: f.TokenEnsurer,
		ActionWaiter: f.ActionWaiter,
		client:       hcapi2.NewMockClient(f.Client),
		config:       f.Config,
	}
}
