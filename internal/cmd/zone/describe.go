package zone

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Zone]{
	ResourceNameSingular: "Zone",
	ShortDescription:     "Describe a Zone",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Zone().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Zone, interface{}, error) {
		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}

		return zone, hcloud.SchemaFromZone(zone), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, zone *hcloud.Zone) error {
		name := util.DisplayZoneName(zone.Name)
		if name != zone.Name {
			name = fmt.Sprintf("%s (IDNA: %s)", name, zone.Name)
		}

		cmd.Printf("ID:\t\t%d\n", zone.ID)
		cmd.Printf("Name:\t\t%s\n", name)
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(zone.Created), humanize.Time(zone.Created))
		cmd.Printf("Mode:\t\t%s\n", zone.Mode)
		cmd.Printf("Status:\t\t%s\n", zone.Status)
		cmd.Printf("TTL:\t\t%d\n", zone.TTL)
		cmd.Printf("Registrar:\t%s\n", zone.Registrar)
		cmd.Printf("Record Count:\t%d\n", zone.RecordCount)
		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(zone.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(zone.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(zone.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Printf("Authoritative Nameservers:\n")
		cmd.Printf("  Assigned:\n")
		if len(zone.AuthoritativeNameservers.Assigned) > 0 {
			for _, srv := range zone.AuthoritativeNameservers.Assigned {
				cmd.Printf("    - %s\n", srv)
			}
		} else {
			cmd.Printf("    No assigned nameservers\n")
		}
		cmd.Printf("  Delegated:\n")
		if len(zone.AuthoritativeNameservers.Delegated) > 0 {
			for _, srv := range zone.AuthoritativeNameservers.Delegated {
				cmd.Printf("    - %s\n", srv)
			}
		} else {
			cmd.Printf("    No delegated nameservers\n")
		}
		cmd.Printf("  Delegation last check:\t%s (%s)\n",
			util.Datetime(zone.AuthoritativeNameservers.DelegationLastCheck),
			humanize.Time(zone.AuthoritativeNameservers.DelegationLastCheck))
		cmd.Printf("  Delegation status:\t\t%s\n", zone.AuthoritativeNameservers.DelegationStatus)

		if zone.Mode == hcloud.ZoneModeSecondary {
			cmd.Printf("Primary nameservers:\n")
			for _, ns := range zone.PrimaryNameservers {
				cmd.Printf("  - Address:\t\t%s\n", ns.Address)
				cmd.Printf("    Port:\t\t%d\n", ns.Port)
				if ns.TSIGAlgorithm != "" {
					cmd.Printf("    TSIG Algorithm:\t%s\n", ns.TSIGAlgorithm)
				}
				if ns.TSIGKey != "" {
					cmd.Printf("    TSIG Key:\t\t%s\n", ns.TSIGKey)
				}
			}
		}
		return nil
	},
	Experimental: experimental.DNS,
}
