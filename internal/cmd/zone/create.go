package zone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd[*hcloud.Zone]{
	BaseCobraCommand: func(_ hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --name <name> [--mode secondary --primary-nameservers <file>]",
			Short:                 "Create a Zone",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("name", "", "Zone name (required)")
		util.MarkFlagRequired(cmd, "name")

		cmd.Flags().String("mode", "primary", "Mode of the Zone (primary, secondary)")
		cmpl.RegisterFlagCompletion(cmd, "mode", cmpl.SuggestCandidates(string(hcloud.ZoneModePrimary), string(hcloud.ZoneModeSecondary)))

		cmd.Flags().Int("ttl", 0, "Default Time To Live (TTL) of the Zone")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmpl.RegisterFlagCompletion(cmd, "enable-protection", cmpl.SuggestCandidates("delete"))

		cmd.Flags().String("primary-nameservers-file", "", "JSON file containing the new primary nameservers. (See 'hcloud zone change-primary-nameservers -h' for help)")
		util.MarkFlagFilename(cmd, "primary-nameservers-file", "json")

		cmd.Flags().String("zonefile", "", "Zone file in BIND (RFC 1034/1035) format (use - to read from stdin)")
		util.MarkFlagFilename(cmd, "zonefile")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, _ []string) (*hcloud.Zone, any, error) {
		name, _ := cmd.Flags().GetString("name")
		mode, _ := cmd.Flags().GetString("mode")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		// Convert name to ascii
		name, err := util.ParseZoneIDOrName(name)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		switch mode {
		case string(hcloud.ZoneModePrimary), string(hcloud.ZoneModeSecondary):
		default:
			return nil, nil, fmt.Errorf("unknown Zone mode: %s", mode)
		}

		protectionOpts, err := ChangeProtectionCmds.GetChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.ZoneCreateOpts{
			Name:   name,
			Mode:   hcloud.ZoneMode(mode),
			Labels: labels,
		}

		if cmd.Flags().Changed("primary-nameservers-file") {
			file, err := cmd.Flags().GetString("primary-nameservers-file")
			if err != nil {
				return nil, nil, err
			}

			nameservers, err := parsePrimaryNameservers(cmd.InOrStdin(), file)
			if err != nil {
				return nil, nil, err
			}

			for _, ns := range nameservers {
				createOpts.PrimaryNameservers = append(createOpts.PrimaryNameservers, hcloud.ZoneCreateOptsPrimaryNameserver{
					Address:       ns.Address,
					Port:          ns.Port,
					TSIGAlgorithm: ns.TSIGAlgorithm,
					TSIGKey:       ns.TSIGKey,
				})
			}
		}

		if cmd.Flags().Changed("ttl") {
			ttl, _ := cmd.Flags().GetInt("ttl")
			createOpts.TTL = &ttl
		}

		if cmd.Flags().Changed("zonefile") {
			if createOpts.Mode == hcloud.ZoneModeSecondary {
				return nil, nil, fmt.Errorf("Zones in secondary mode can not be created from a zone file")
			}

			zonefile, _ := cmd.Flags().GetString("zonefile")
			createOpts.Zonefile, err = readZonefile(cmd.InOrStdin(), zonefile)
			if err != nil {
				return nil, nil, err
			}
		}

		result, _, err := s.Client().Zone().Create(cmd.Context(), createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Zone %s created\n", result.Zone.Name)

		if protectionOpts.Delete != nil {
			if err := ChangeProtectionCmds.ChangeProtection(s, cmd, result.Zone, true, protectionOpts); err != nil {
				return nil, nil, err
			}
		}

		// Assigned authoritative nameserver is only set after the action completed. Need to reload the zone
		zone, _, err := s.Client().Zone().GetByID(cmd.Context(), result.Zone.ID)
		if err != nil {
			return nil, nil, err
		}

		return zone, util.Wrap("zone", hcloud.SchemaFromZone(zone)), nil
	},
}
