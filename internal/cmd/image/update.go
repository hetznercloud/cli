package image

import (
	"context"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var updateCmd = base.UpdateCmd{
	ResourceNameSingular: "Image",
	ShortDescription:     "Update a Image",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Image().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Image().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("description", "", "Image description")
		cmd.Flags().String("type", "", "Image type")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("backup", "snapshot"))
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		image := resource.(*hcloud.Image)
		updOpts := hcloud.ImageUpdateOpts{
			Description: hcloud.String(flags["description"].String()),
			Type:        hcloud.ImageType(flags["type"].String()),
		}
		_, _, err := client.Image().Update(ctx, image, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
