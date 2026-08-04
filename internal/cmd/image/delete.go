package image

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.Image]{
	ResourceNameSingular: "Image",
	ResourceNamePlural:   "Images",
	ShortDescription:     "Delete an Image",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Image().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Image, *hcloud.Response, error) {
		id, err := strconv.ParseInt(idOrName, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid snapshot or backup ID %q", idOrName)
		}
		return s.Client().Image().GetByID(cmd.Context(), id)
	},
	Delete: func(s state.State, cmd *cobra.Command, image *hcloud.Image) ([]*hcloud.Action, error) {
		_, err := s.Client().Image().Delete(cmd.Context(), image)
		return nil, err
	},
}
