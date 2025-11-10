package zone

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ChangePrimaryNameserversCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "change-primary-nameservers --primary-nameservers-file <file> <zone>",
			Short: "Changes the primary nameservers of a secondary Zone",
			Long: `Changes the primary nameservers of a secondary Zone.

Input file has to be in JSON format. You can find the schema at https://docs.hetzner.cloud/reference/cloud#zone-actions-change-a-zone-primary-nameservers

Example file content:

[
  {
    "address": "203.0.113.10"
  },
  {
    "address": "203.0.113.11",
    "port": 5353
  },
  {
    "address": "203.0.113.12",
    "tsig_algorithm": "hmac-sha256",
    "tsig_key": "example-key"
  }
]`,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Zone().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("primary-nameservers-file", "", "JSON file containing the new primary nameservers. (use - to read from stdin)")
		_ = cmd.MarkFlagRequired("primary-nameservers-file")
		_ = cmd.MarkFlagFilename("primary-nameservers-file", "json")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, idOrName)
		if err != nil {
			return err
		}
		if zone == nil {
			return fmt.Errorf("Zone not found: %s", idOrName)
		}

		var opts hcloud.ZoneChangePrimaryNameserversOpts

		file, _ := cmd.Flags().GetString("primary-nameservers-file")
		opts.PrimaryNameservers, err = parsePrimaryNameservers(file)
		if err != nil {
			return err
		}

		action, _, err := s.Client().Zone().ChangePrimaryNameservers(s, zone, opts)
		if err != nil {
			return err
		}
		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Primary nameservers for Zone %s updated\n", zone.Name)
		return nil
	},
}

func parsePrimaryNameservers(path string) ([]hcloud.ZoneChangePrimaryNameserversOptsPrimaryNameserver, error) {
	var (
		data []byte
		err  error
	)
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}

	var nameserverSchemas []schema.ZonePrimaryNameserver
	err = json.Unmarshal(data, &nameserverSchemas)
	if err != nil {
		return nil, err
	}

	var nameservers []hcloud.ZoneChangePrimaryNameserversOptsPrimaryNameserver
	for _, ns := range nameserverSchemas {
		nameservers = append(nameservers, hcloud.ZoneChangePrimaryNameserversOptsPrimaryNameserver{
			Address:       ns.Address,
			Port:          ns.Port,
			TSIGAlgorithm: hcloud.ZoneTSIGAlgorithm(ns.TSIGAlgorithm),
			TSIGKey:       ns.TSIGKey,
		})
	}

	return nameservers, err
}
