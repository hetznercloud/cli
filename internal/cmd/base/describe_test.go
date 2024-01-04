package base_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
)

var fakeDescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "Fake resource",

	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		cmd.Println("Fetching fake resource")
		commandCalled = true

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, util.Wrap("resource", resource), nil
	},

	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		rsc := resource.(*fakeResource)
		cmd.Printf("ID: %d\n", rsc.ID)
		cmd.Printf("Name: %s\n", rsc.Name)
		return nil
	},

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return nil
	},
}

func TestDescribe(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Equal(t, `Fetching fake resource
ID: 123
Name: test
`, out)
	assert.Empty(t, errOut)
}

func TestDescribeJSON(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123", "-o=json"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Equal(t, "Fetching fake resource\n", errOut)
}

func TestDescribeYAML(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123", "-o=yaml"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Equal(t, "Fetching fake resource\n", errOut)
}

func TestDescribeQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, errOut)
}

func TestDescribeJSONQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123", "-o=json", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Empty(t, errOut)
}

func TestDescribeYAMLQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeDescribeCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"describe", "123", "-o=yaml", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Empty(t, errOut)
}
