package zone

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ImportZonefileCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "import-zonefile --zonefile <file> <zone>",
			Short:                 "Imports a zone file, replacing all Zone RRSets",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Zone().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("zonefile", "", "Zone file in BIND (RFC 1034/1035) format (use - to read from stdin)")
		_ = cmd.MarkFlagRequired("zonefile")
		_ = cmd.MarkFlagFilename("zonefile")

		output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())

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

		opts := hcloud.ZoneImportZonefileOpts{}

		zonefile, _ := cmd.Flags().GetString("zonefile")
		opts.Zonefile, err = readZonefile(zonefile)
		if err != nil {
			return err
		}

		action, _, err := s.Client().Zone().ImportZonefile(s, zone, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Zone file for Zone %s imported\n", zone.Name)

		return nil
	},
	Experimental: experimental.DNS,
}

func readZonefile(zonefile string) (string, error) {
	var data []byte
	var err error

	if zonefile == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(zonefile)
	}
	if err != nil {
		return "", err
	}

	return string(data), nil
}
