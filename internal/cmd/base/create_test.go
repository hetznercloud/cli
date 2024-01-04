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

type fakeResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var commandCalled bool

var fakeCreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use: "create",
		}
	},
	Run: func(s state.State, cmd *cobra.Command, strings []string) (any, any, error) {
		cmd.Println("Creating fake resource")
		commandCalled = true

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, util.Wrap("resource", resource), nil
	},
}

func TestCreate(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Equal(t, "Creating fake resource\n", out)
	assert.Empty(t, errOut)
}

func TestCreateJSON(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create", "-o=json"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Equal(t, "Creating fake resource\n", errOut)
}

func TestCreateYAML(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create", "-o=yaml"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Equal(t, "Creating fake resource\n", errOut)
}

func TestCreateQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, errOut)
}

func TestCreateJSONQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create", "-o=json", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Empty(t, errOut)
}

func TestCreateYAMLQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeCreateCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"create", "-o=yaml", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `{"resource": {"id": 123, "name": "test"}}`, out)
	assert.Empty(t, errOut)
}
