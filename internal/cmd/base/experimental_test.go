package base_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
)

type fakeExperimentalCmd struct{}

func (fakeExperimentalCmd) CobraCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "experimental",
		Short:   "My experimental command",
		Long:    "This is an experimental command that may change in the future.",
		PreRunE: s.EnsureToken,
	}

	cmd.Run = func(cmd *cobra.Command, _ []string) {
		cmd.Println("Hello world")
	}

	return base.Experimental(s, cmd, "Example Product", "https://example.com")
}

func TestExperimental(t *testing.T) {
	testutil.TestCommand(t, fakeExperimentalCmd{}, map[string]testutil.TestCase{
		"default": {
			Args:      []string{"experimental"},
			ExpOut:    "Hello world\n",
			ExpErrOut: "Warning: This command is experimental and may change in the future. Use --experimental to suppress this warning.\n",
		},
		"experimental": {
			Args:   []string{"experimental", "--experimental"},
			ExpOut: "Hello world\n",
		},
	})

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := fakeExperimentalCmd{}.CobraCommand(fx.State())
	assert.Equal(t, "[experimental] My experimental command", cmd.Short)
	assert.Equal(t, `This is an experimental command that may change in the future.

Experimental: Example Product is experimental, breaking changes may occur within minor releases. See https://example.com for more details.`, cmd.Long)
}
