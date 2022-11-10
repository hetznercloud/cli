package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"
	hcapi2_mock "github.com/hetznercloud/cli/internal/hcapi2/mock"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
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
func (f *Fixture) Run(cmd *cobra.Command, args []string) (string, error) {
	cmd.SetArgs(args)
	return CaptureStdout(cmd.Execute)
}

// Finish must be called after the test is finished, preferably via `defer` directly after
// creating the Fixture.
func (f *Fixture) Finish() {
	f.MockController.Finish()
}
