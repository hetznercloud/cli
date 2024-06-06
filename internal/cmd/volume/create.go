package volume

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/actionutils"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create [options] --name <name> --size <size>",
			Short:                 "Create a volume",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("name", "", "Volume name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().String("server", "", "Server (ID or name)")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().String("location", "", "Location (ID or name)")
		cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().Int("size", 0, "Size (GB) (required)")
		cmd.MarkFlagRequired("size")

		cmd.Flags().Bool("automount", false, "Automount volume after attach (server must be provided)")

		cmd.Flags().String("format", "", "Format volume after creation (ext4 or xfs)")
		cmd.RegisterFlagCompletionFunc("format", cmpl.SuggestCandidates("ext4", "xfs"))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringSlice("enable-protection", []string{}, "Enable protection (delete) (default: none)")
		cmd.RegisterFlagCompletionFunc("enable-protection", cmpl.SuggestCandidates("delete"))

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (any, any, error) {
		name, _ := cmd.Flags().GetString("name")
		serverIDOrName, _ := cmd.Flags().GetString("server")
		size, _ := cmd.Flags().GetInt("size")
		location, _ := cmd.Flags().GetString("location")
		automount, _ := cmd.Flags().GetBool("automount")
		format, _ := cmd.Flags().GetString("format")
		labels, _ := cmd.Flags().GetStringToString("label")
		protection, _ := cmd.Flags().GetStringSlice("enable-protection")

		protectionOpts, err := getChangeProtectionOpts(true, protection)
		if err != nil {
			return nil, nil, err
		}

		createOpts := hcloud.VolumeCreateOpts{
			Name:   name,
			Size:   size,
			Labels: labels,
		}

		if location != "" {
			id, err := strconv.ParseInt(location, 10, 64)
			if err == nil {
				createOpts.Location = &hcloud.Location{ID: id}
			} else {
				createOpts.Location = &hcloud.Location{Name: location}
			}
		}
		if serverIDOrName != "" {
			server, _, err := s.Client().Server().Get(s, serverIDOrName)
			if err != nil {
				return nil, nil, err
			}
			if server == nil {
				return nil, nil, fmt.Errorf("server not found: %s", serverIDOrName)
			}
			createOpts.Server = server
		}
		if automount {
			createOpts.Automount = &automount
		}
		if format != "" {
			createOpts.Format = &format
		}

		result, _, err := s.Client().Volume().Create(s, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(cmd, s, actionutils.AppendNextActions(result.Action, result.NextActions)...); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Volume %d created\n", result.Volume.ID)

		if err := changeProtection(s, cmd, result.Volume, true, protectionOpts); err != nil {
			return nil, nil, err
		}

		return result.Volume, util.Wrap("volume", hcloud.SchemaFromVolume(result.Volume)), nil
	},
}
