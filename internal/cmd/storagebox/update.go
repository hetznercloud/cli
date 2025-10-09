package storagebox

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Storage Box",
	ShortDescription:     "Update a Storage Box",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.StorageBox().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (any, *hcloud.Response, error) {
		return s.Client().StorageBox().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Storage Box name")
	},
	Update: func(s state.State, _ *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		storageBox := resource.(*hcloud.StorageBox)
		opts := hcloud.StorageBoxUpdateOpts{
			Name: hcloud.Ptr(flags["name"].String()),
		}
		_, _, err := s.Client().StorageBox().Update(s, storageBox, opts)
		if err != nil {
			return err
		}
		return nil
	},
}
