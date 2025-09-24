package base_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
)

var ExperimentalProduct = base.ExperimentalWrapper("Product name", "experimental", "https://docs.hetzner.cloud/changelog#new-product")

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

	return ExperimentalProduct(s, cmd)
}

func TestExperimental(t *testing.T) {
	testutil.TestCommand(t, fakeExperimentalCmd{}, map[string]testutil.TestCase{
		"default": {
			Args:      []string{"experimental"},
			ExpOut:    "Hello world\n",
			ExpErrOut: "Warning: Product name is experimental and may change in the future. Use --no-experimental-warnings to suppress this warning.\n",
		},
		"experimental": {
			Args:   []string{"experimental", "--no-experimental-warnings"},
			ExpOut: "Hello world\n",
		},
	})

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := fakeExperimentalCmd{}.CobraCommand(fx.State())
	assert.Equal(t, "[experimental] My experimental command", cmd.Short)
	assert.Equal(t, `This is an experimental command that may change in the future.

Experimental: Product name is experimental, breaking changes may occur within minor releases.
See https://docs.hetzner.cloud/changelog#new-product for more details.
`, cmd.Long)
}
