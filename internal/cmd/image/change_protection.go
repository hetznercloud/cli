package image

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.Image, hcloud.ImageChangeProtectionOpts]{
	ResourceNameSingular:    "Image",
	ShortEnableDescription:  "Enable resource protection for an Image",
	ShortDisableDescription: "Disable resource protection for an Image",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.Image().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.ImageChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.ImageChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(_ state.State, _ *cobra.Command, idStr string) (*hcloud.Image, *hcloud.Response, error) {
		imageID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, nil, errors.New("invalid Image ID")
		}
		return &hcloud.Image{ID: imageID}, nil, nil
	},

	ChangeProtectionFunction: func(s state.State, image *hcloud.Image, opts hcloud.ImageChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Image().ChangeProtection(s, image, opts)
	},

	IDOrName: func(image *hcloud.Image) string {
		return fmt.Sprint(image.ID)
	},
}
