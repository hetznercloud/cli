package base_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var fakeDeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Fake resource",
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		cmd.Println("Deleting fake resource")
		commandCalled = true
		return nil
	},

	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		cmd.Println("Fetching fake resource")

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, nil, nil
	},

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return nil
	},
}

func TestDelete(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDeleteCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"delete", "123"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Equal(t, `Fetching fake resource
Deleting fake resource
Fake resource 123 deleted
`, out)
	assert.Empty(t, errOut)
}

func TestDeleteQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDeleteCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"delete", "123", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, errOut)
}
