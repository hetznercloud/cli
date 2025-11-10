package zone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Zone",
	ResourceNamePlural:   "Zones",
	ShortDescription:     "Delete a Zone",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Zone().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		return s.Client().Zone().Get(s, idOrName)
	},
	Delete: func(s state.State, _ *cobra.Command, resource interface{}) (*hcloud.Action, error) {
		zone := resource.(*hcloud.Zone)
		res, _, err := s.Client().Zone().Delete(s, zone)
		if err != nil {
			return nil, err
		}
		return res.Action, nil
	},
}
