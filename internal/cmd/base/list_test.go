package base_test

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var fakeListCmd = &base.ListCmd{
	ResourceNamePlural: "Fake resources",

	Schema: func(i []interface{}) interface{} {
		return i
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
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

	Fetch: func(_ state.State, _ *pflag.FlagSet, _ hcloud.ListOpts, sort []string) ([]interface{}, error) {
		resources := []*fakeResource{
			{
				ID:   456,
				Name: "test2",
			},
			{
				ID:   123,
				Name: "test",
			},
			{
				ID:   789,
				Name: "test3",
			},
		}
		if len(sort) > 0 {
			switch sort[0] {
			case "id:asc":
				slices.SortFunc(resources, func(a, b *fakeResource) int {
					return cmp.Compare(a.ID, b.ID)
				})
			case "id:desc":
				slices.SortFunc(resources, func(a, b *fakeResource) int {
					return cmp.Compare(b.ID, a.ID)
				})
			case "name:asc":
				slices.SortFunc(resources, func(a, b *fakeResource) int {
					return cmp.Compare(a.Name, b.Name)
				})
			case "name:desc":
				slices.SortFunc(resources, func(a, b *fakeResource) int {
					return cmp.Compare(b.Name, a.Name)
				})
			}
		}
		return util.ToAnySlice(resources), nil
	},
}

func TestList(t *testing.T) {
	sortOpt, cleanup := config.NewTestOption(
		"sort.fakeresource",
		"",
		[]string{"id:asc"},
		(config.DefaultPreferenceFlags&^config.OptionFlagPFlag)|config.OptionFlagSlice,
		nil,
	)
	defer cleanup()

	fakeListCmd.SortOption = sortOpt

	const resourceSchema = `[{"id": 123, "name": "test"}, {"id": 456, "name": "test2"}, {"id": 789, "name": "test3"}]`
	testutil.TestCommand(t, fakeListCmd, map[string]testutil.TestCase{
		"no flags": {
			Args:   []string{"list"},
			ExpOut: "ID    NAME\n123   test\n456   test2\n789   test3\n",
		},
		"json": {
			Args:       []string{"list", "-o=json"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
		},
		"yaml": {
			Args:       []string{"list", "-o=yaml"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
		},
		"quiet": {
			Args:   []string{"list", "--quiet"},
			ExpOut: "ID    NAME\n123   test\n456   test2\n789   test3\n",
		},
		"json quiet": {
			Args:       []string{"list", "-o=json", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeJSON,
		},
		"yaml quiet": {
			Args:       []string{"list", "-o=yaml", "--quiet"},
			ExpOut:     resourceSchema,
			ExpOutType: testutil.DataTypeYAML,
		},
		"sort": {
			Args:       []string{"list", "--sort", "id:desc", "-o=json"},
			ExpOut:     `[{"id": 789, "name": "test3"}, {"id": 456, "name": "test2"}, {"id": 123, "name": "test"}]`,
			ExpOutType: testutil.DataTypeJSON,
		},
		"no sort": {
			Args:       []string{"list", "--sort=", "-o=json"},
			ExpOut:     `[{"id": 456, "name": "test2"}, {"id": 123, "name": "test"}, {"id": 789, "name": "test3"}]`,
			ExpOutType: testutil.DataTypeJSON,
		},
		"sort with option": {
			Args: []string{"list", "-o=json"},
			PreRun: func(t *testing.T, fx *testutil.Fixture) {
				sortOpt.Override(fx.Config, []string{"id:desc"})
				t.Cleanup(func() {
					sortOpt.Override(fx.Config, nil)
				})
			},
			ExpOut:     `[{"id": 789, "name": "test3"}, {"id": 456, "name": "test2"}, {"id": 123, "name": "test"}]`,
			ExpOutType: testutil.DataTypeJSON,
		},
	})
}
