package zone

import (
	"fmt"
	"io"

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
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, zone *hcloud.Zone) error {

		name := util.DisplayZoneName(zone.Name)
		if name != zone.Name {
			name = fmt.Sprintf("%s (IDNA: %s)", name, zone.Name)
		}

		_, _ = fmt.Fprintf(out, "ID:\t%d\n", zone.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(zone.Created), humanize.Time(zone.Created))
		_, _ = fmt.Fprintf(out, "Mode:\t%s\n", zone.Mode)
		_, _ = fmt.Fprintf(out, "Status:\t%s\n", zone.Status)
		_, _ = fmt.Fprintf(out, "TTL:\t%d\n", zone.TTL)
		_, _ = fmt.Fprintf(out, "Registrar:\t%s\n", zone.Registrar)
		_, _ = fmt.Fprintf(out, "Record Count:\t%d\n", zone.RecordCount)
		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%s\n", util.YesNo(zone.Protection.Delete))

		util.DescribeLabels(out, zone.Labels, "")

		_, _ = fmt.Fprintf(out, "\n")

		_, _ = fmt.Fprintf(out, "Authoritative Nameservers:\n")
		_, _ = fmt.Fprintf(out, "  Assigned:\n")
		if len(zone.AuthoritativeNameservers.Assigned) > 0 {
			for _, srv := range zone.AuthoritativeNameservers.Assigned {
				_, _ = fmt.Fprintf(out, "    - %s\n", srv)
			}
		} else {
			_, _ = fmt.Fprintf(out, "    No assigned nameservers\n")
		}

		_, _ = fmt.Fprintf(out, "  Delegated:\n")
		if len(zone.AuthoritativeNameservers.Delegated) > 0 {
			for _, srv := range zone.AuthoritativeNameservers.Delegated {
				_, _ = fmt.Fprintf(out, "    - %s\n", srv)
			}
		} else {
			_, _ = fmt.Fprintf(out, "    No delegated nameservers\n")
		}
		_, _ = fmt.Fprintf(out, "  Delegation last check:\t%s (%s)\n",
			util.Datetime(zone.AuthoritativeNameservers.DelegationLastCheck),
			humanize.Time(zone.AuthoritativeNameservers.DelegationLastCheck))
		_, _ = fmt.Fprintf(out, "  Delegation status:\t%s\n", zone.AuthoritativeNameservers.DelegationStatus)

		if zone.Mode == hcloud.ZoneModeSecondary {
			_, _ = fmt.Fprintf(out, "Primary nameservers:\t\n")
			for _, ns := range zone.PrimaryNameservers {
				_, _ = fmt.Fprintf(out, "  - Address:\t%s\n", ns.Address)
				_, _ = fmt.Fprintf(out, "    Port:\t%d\n", ns.Port)
				if ns.TSIGAlgorithm != "" {
					_, _ = fmt.Fprintf(out, "    TSIG Algorithm:\t%s\n", ns.TSIGAlgorithm)
				}
				if ns.TSIGKey != "" {
					_, _ = fmt.Fprintf(out, "    TSIG Key:\t%s\n", ns.TSIGKey)
				}
			}
		}
		return nil
	},
	Experimental: experimental.DNS,
}
