package zone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var ExportZonefileCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "export-zonefile [options] <zone>",
			Short:                 "Returns a generated Zone file in BIND (RFC 1034/1035) format",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Zone().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		output.AddFlag(cmd, output.OptionJSON(), output.OptionYAML())
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		outputFlags := output.FlagsForCommand(cmd)

		zone, _, err := s.Client().Zone().Get(s, idOrName)
		if err != nil {
			return err
		}
		if zone == nil {
			return fmt.Errorf("Zone not found: %s", idOrName)
		}

		res, _, err := s.Client().Zone().ExportZonefile(s, zone)
		if err != nil {
			return err
		}

		schema := util.Wrap("zonefile", res.Zonefile)

		switch {
		case outputFlags.IsSet("json"):
			return util.DescribeJSON(cmd.OutOrStdout(), schema)
		case outputFlags.IsSet("yaml"):
			return util.DescribeYAML(cmd.OutOrStdout(), schema)
		default:
			cmd.Print(res.Zonefile)
		}
		return nil
	},
}
