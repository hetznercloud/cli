package base_test

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var mu = sync.Mutex{}

var fakeDeleteCmd = &base.DeleteCmd[*fakeResource]{
	ResourceNameSingular: "Fake resource",
	ResourceNamePlural:   "Fake resources",
	Delete: func(_ state.State, cmd *cobra.Command, _ *fakeResource) ([]*hcloud.Action, error) {
		defer mu.Unlock()
		cmd.Println("Deleting fake resource")
		return nil, nil
	},

	Fetch: func(_ state.State, cmd *cobra.Command, idOrName string) (*fakeResource, *hcloud.Response, error) {
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

	NameSuggestions: func(hcapi2.Client) hcapi2.CompletionFunc {
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

func TestDeleteReportsActionFailureForOwningResource(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := (&base.DeleteCmd[*fakeResource]{
		ResourceNameSingular: "Fake resource",
		ResourceNamePlural:   "Fake resources",
		NameSuggestions:      func(hcapi2.Client) hcapi2.CompletionFunc { return nil },
		Fetch: func(_ state.State, _ *cobra.Command, idOrName string) (*fakeResource, *hcloud.Response, error) {
			id, err := strconv.Atoi(idOrName)
			if err != nil {
				return nil, nil, err
			}
			return &fakeResource{ID: id, Name: idOrName}, nil, nil
		},
		Delete: func(_ state.State, _ *cobra.Command, resource *fakeResource) ([]*hcloud.Action, error) {
			return []*hcloud.Action{{ID: int64(resource.ID)}}, nil
		},
	}).CobraCommand(fx.State())
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	fx.ExpectEnsureToken()
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 1}, &hcloud.Action{ID: 2}).
		Return(&state.ActionWaitError{Failures: []state.ActionFailure{{ActionID: 2, Err: errors.New("action failed")}}})

	out, errOut, err := fx.Run(cmd, []string{"1", "2"})
	require.EqualError(t, err, "Fake resource 2: action 2 failed: action failed")
	assert.Equal(t, "Fake resource 1 deleted\n", out)
	assert.Empty(t, errOut)
}
