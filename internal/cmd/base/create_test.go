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

type fakeResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var fakeCreateCmd = &base.CreateCmd[*fakeResource]{
	BaseCobraCommand: func(hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use: "create",
		}
	},
	Run: func(_ state.State, cmd *cobra.Command, _ []string) (*fakeResource, any, error) {
		cmd.Println("Creating fake resource")

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, util.Wrap("resource", resource), nil
	},
}

func TestCreate(t *testing.T) {
	const resourceSchema = `{"resource": {"id": 123, "name": "test"}}`
	testutil.TestCommand(t, fakeCreateCmd, map[string]testutil.TestCase{
		"no flags": {
			Args:   []string{"create"},
			ExpOut: "Creating fake resource\n",
		},
		"json": {
			Args:       []string{"create", "-o=json"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
			ExpErrOut:  "Creating fake resource\n",
		},
		"yaml": {
			Args:       []string{"create", "-o=yaml"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
			ExpErrOut:  "Creating fake resource\n",
		},
		"quiet": {
			Args: []string{"create", "--quiet"},
		},
		"json quiet": {
			Args:       []string{"create", "-o=json", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
		},
		"yaml quiet": {
			Args:       []string{"create", "-o=yaml", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
		},
	})
}
