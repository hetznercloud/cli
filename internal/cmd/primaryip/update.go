package primaryip

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Primary IP",
	ShortDescription:     "Update a Primary IP",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().PrimaryIP().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Primary IP name")
		cmd.Flags().Bool("auto-delete", false, "Delete this Primary IP when the resource it is assigned to is deleted")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		primaryIP := resource.(*hcloud.PrimaryIP)
		updOpts := hcloud.PrimaryIPUpdateOpts{
			Name: flags["name"].String(),
		}

		if cmd.Flags().Changed("auto-delete") {
			autoDelete, _ := cmd.Flags().GetBool("auto-delete")
			updOpts.AutoDelete = hcloud.Ptr(autoDelete)
		}

		_, _, err := s.Client().PrimaryIP().Update(s, primaryIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
