package testutil

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	hcapi2_mock "github.com/hetznercloud/cli/internal/hcapi2/mock"
	"github.com/hetznercloud/cli/internal/state"
)

// Fixture provides affordances for testing CLI commands.
type Fixture struct {
	MockController *gomock.Controller
	Client         *hcapi2_mock.MockClient
	ActionWaiter   *state.MockActionWaiter
	TokenEnsurer   *state.MockTokenEnsurer
}

// NewFixture creates a new Fixture.
func NewFixture(t *testing.T) *Fixture {
	ctrl := gomock.NewController(t)

	return &Fixture{
		MockController: ctrl,
		Client:         hcapi2_mock.NewMockClient(ctrl),
		ActionWaiter:   state.NewMockActionWaiter(ctrl),
		TokenEnsurer:   state.NewMockTokenEnsurer(ctrl),
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
	state.TokenEnsurer
	state.ActionWaiter
	context.Context
	hcapi2.Client
}

func (fixtureState) WriteConfig() error {
	return errors.New("not implemented")
}

func (fixtureState) Config() *state.Config {
	return nil
}

// State returns a state.State implementation for testing purposes.
func (f *Fixture) State() state.State {
	return fixtureState{
		TokenEnsurer: f.TokenEnsurer,
		ActionWaiter: f.ActionWaiter,
		Context:      context.Background(),
		Client:       f.Client,
	}
}
