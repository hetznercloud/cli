package image

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Image",
	ShortDescription:     "Update an image",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Image().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().Image().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("description", "", "Image description")
		cmd.Flags().String("type", "", "Image type")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates("snapshot"))
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		image := resource.(*hcloud.Image)
		updOpts := hcloud.ImageUpdateOpts{
			Description: hcloud.String(flags["description"].String()),
			Type:        hcloud.ImageType(flags["type"].String()),
		}
		_, _, err := s.Client().Image().Update(s, image, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
