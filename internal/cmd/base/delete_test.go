package base_test

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var fakeDeleteCmd = &base.DeleteCmd{
	ResourceNameSingular: "Fake resource",
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		cmd.Println("Deleting fake resource")
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
	testutil.TestCommand(t, fakeDeleteCmd, map[string]testutil.TestCase{
		"no flags": {
			Args:   []string{"delete", "123"},
			ExpOut: "Fetching fake resource\nDeleting fake resource\nFake resource 123 deleted\n",
		},
		"quiet": {
			Args: []string{"delete", "123", "--quiet"},
		},
	})
}
