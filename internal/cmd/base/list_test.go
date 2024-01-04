package base_test

import (
	"fmt"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var fakeListCmd = base.ListCmd{
	ResourceNamePlural: "Fake resources",

	Schema: func(i []interface{}) interface{} {
		return i
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.Firewall{}).
			AddFieldFn("id", func(obj interface{}) string {
				rsc := obj.(*fakeResource)
				return fmt.Sprintf("%d", rsc.ID)
			}).
			AddFieldFn("name", func(obj interface{}) string {
				rsc := obj.(*fakeResource)
				return rsc.Name
			})
	},

	DefaultColumns: []string{"id", "name"},

	Fetch: func(s state.State, set *pflag.FlagSet, opts hcloud.ListOpts, strings []string) ([]interface{}, error) {
		commandCalled = true
		return []interface{}{
			&fakeResource{
				ID:   123,
				Name: "test",
			},
			&fakeResource{
				ID:   321,
				Name: "test2",
			},
			&fakeResource{
				ID:   42,
				Name: "test3",
			},
		}, nil
	},
}

func TestList(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Equal(t, "ID    NAME\n123   test\n321   test2\n42    test3\n", out)
	assert.Empty(t, errOut)
}

func TestListJSON(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list", "-o=json"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `
[
  {
    "id": 123,
    "name": "test"
  },
  {
    "id": 321,
    "name": "test2"
  },
  {
    "id": 42,
    "name": "test3"
  }
]`, out)
	assert.Empty(t, errOut)
}

func TestListYAML(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list", "-o=json"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `
[
  {
    "id": 123,
    "name": "test"
  },
  {
    "id": 321,
    "name": "test2"
  },
  {
    "id": 42,
    "name": "test3"
  }
]`, out)
	assert.Empty(t, errOut)
}

func TestListQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.Equal(t, "ID    NAME\n123   test\n321   test2\n42    test3\n", out)
	assert.Empty(t, errOut)
}

func TestListJSONQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list", "-o=json", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.JSONEq(t, `
[
  {
    "id": 123,
    "name": "test"
  },
  {
    "id": 321,
    "name": "test2"
  },
  {
    "id": 42,
    "name": "test3"
  }
]`, out)
	assert.Empty(t, errOut)
}

func TestListYAMLQuiet(t *testing.T) {
	commandCalled = false

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.AddCommand(fakeListCmd.CobraCommand(fx.State()))

	out, errOut, err := fx.Run(cmd, []string{"list", "-o=json", "--quiet"})

	assert.Equal(t, true, commandCalled)
	assert.NoError(t, err)
	assert.YAMLEq(t, `
[
  {
    "id": 123,
    "name": "test"
  },
  {
    "id": 321,
    "name": "test2"
  },
  {
    "id": 42,
    "name": "test3"
  }
]`, out)
	assert.Empty(t, errOut)
}
