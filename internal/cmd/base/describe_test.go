package base_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
)

var fakeDescribeCmd = &base.DescribeCmd{
	ResourceNameSingular: "Fake resource",

	Fetch: func(_ state.State, cmd *cobra.Command, _ string) (interface{}, interface{}, error) {
		cmd.Println("Fetching fake resource")

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, util.Wrap("resource", resource), nil
	},

	PrintText: func(_ state.State, cmd *cobra.Command, resource interface{}) error {
		rsc := resource.(*fakeResource)
		cmd.Printf("ID: %d\n", rsc.ID)
		cmd.Printf("Name: %s\n", rsc.Name)
		return nil
	},

	NameSuggestions: func(hcapi2.Client) func() []string {
		return nil
	},
}

func TestDescribe(t *testing.T) {
	const resourceSchema = `{"resource": {"id": 123, "name": "test"}}`
	testutil.TestCommand(t, fakeDescribeCmd, map[string]testutil.TestCase{
		"no flags": {
			Args: []string{"describe", "123"},
			ExpOut: `Fetching fake resource
ID: 123
Name: test
`,
		},
		"json": {
			Args:       []string{"describe", "123", "-o=json"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
			ExpErrOut:  "Fetching fake resource\n",
		},
		"yaml": {
			Args:       []string{"describe", "123", "-o=yaml"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
			ExpErrOut:  "Fetching fake resource\n",
		},
		"quiet": {
			Args: []string{"describe", "123", "--quiet"},
		},
		"json quiet": {
			Args:       []string{"describe", "123", "-o=json", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
		},
		"yaml quiet": {
			Args:       []string{"describe", "123", "-o=yaml", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
		},
	})
}
