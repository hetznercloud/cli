package zone

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.Zone]{
	ResourceNameSingular:   "Zone",
	ShortDescriptionAdd:    "Add a label to a Zone",
	ShortDescriptionRemove: "Remove a label from a Zone",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Zone().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Zone().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.Zone, error) {
		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if zone == nil {
			return nil, fmt.Errorf("Zone not found: %s", idOrName)
		}

		return zone, nil
	},
	SetLabels: func(s state.State, zone *hcloud.Zone, labels map[string]string) error {
		opts := hcloud.ZoneUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Zone().Update(s, zone, opts)
		return err
	},
	GetLabels: func(zone *hcloud.Zone) map[string]string {
		return zone.Labels
	},
	GetIDOrName: func(zone *hcloud.Zone) string {
		return zone.Name
	},
}
