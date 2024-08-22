package base_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var mu = sync.Mutex{}

var fakeDeleteCmd = &base.DeleteCmd{
	ResourceNameSingular: "Fake resource",
	ResourceNamePlural:   "Fake resources",
	Delete: func(_ state.State, cmd *cobra.Command, _ interface{}) (*hcloud.Action, error) {
		defer mu.Unlock()
		cmd.Println("Deleting fake resource")
		return nil, nil
	},

	Fetch: func(_ state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		mu.Lock()
		cmd.Println("Fetching fake resource")

		if idOrName == "fail" {
			mu.Unlock()
			return nil, nil, errors.New("this is an error")
		}

		resource := &fakeResource{
			ID:   123,
			Name: "test",
		}

		return resource, nil, nil
	},

	NameSuggestions: func(hcapi2.Client) func() []string {
		return nil
	},
}

func TestDelete(t *testing.T) {
	testutil.TestCommand(t, fakeDeleteCmd, map[string]testutil.TestCase{
		"no flags": {
			Args:   []string{"delete", "123"},
			ExpOut: "Fetching fake resource\nDeleting fake resource\nFake resource 123 deleted\n",
		},
		"no flags multiple": {
			Args: []string{"delete", "123", "456", "789"},
			ExpOut: "Fetching fake resource\nDeleting fake resource\nFetching fake resource\nDeleting fake resource\n" +
				"Fetching fake resource\nDeleting fake resource\nFake resources 123, 456, 789 deleted\n",
		},
		"error": {
			Args:   []string{"delete", "fail"},
			ExpOut: "Fetching fake resource\n",
			ExpErr: "this is an error",
		},
		"error multiple": {
			Args:   []string{"delete", "123", "fail", "789"},
			ExpOut: "Fetching fake resource\nDeleting fake resource\nFetching fake resource\nFetching fake resource\nDeleting fake resource\nFake resources 123, 789 deleted\n",
			ExpErr: "this is an error",
		},
		"quiet": {
			Args: []string{"delete", "123", "--quiet"},
		},
		"quiet multiple": {
			Args: []string{"delete", "123", "456", "789", "--quiet"},
		},
	})
}
