package base_test

import (
	"fmt"
	"testing"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var fakeListCmd = &base.ListCmd{
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
	const resourceSchema = `[{"id": 123, "name": "test"}, {"id": 321, "name": "test2"}, {"id": 42, "name": "test3"}]`
	testutil.TestCommand(t, fakeListCmd, map[string]testutil.TestCase{
		"no flags": {
			Args:   []string{"list"},
			ExpOut: "ID    NAME\n123   test\n321   test2\n42    test3\n",
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
			ExpOut: "ID    NAME\n123   test\n321   test2\n42    test3\n",
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
	})
}
