package volume

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var CreateCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "create FLAGS",
			Short:                 "Create a volume",
			Args:                  cobra.NoArgs,
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

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		serverIDOrName, _ := cmd.Flags().GetString("server")
		size, _ := cmd.Flags().GetInt("size")
		location, _ := cmd.Flags().GetString("location")
		automount, _ := cmd.Flags().GetBool("automount")
		format, _ := cmd.Flags().GetString("format")
		labels, _ := cmd.Flags().GetStringToString("label")

		opts := hcloud.VolumeCreateOpts{
			Name:   name,
			Size:   size,
			Labels: labels,
		}

		if location != "" {
			id, err := strconv.ParseInt(location, 10, 64)
			if err == nil {
				opts.Location = &hcloud.Location{ID: id}
			} else {
				opts.Location = &hcloud.Location{Name: location}
			}
		}
		if serverIDOrName != "" {
			server, _, err := client.Server().Get(ctx, serverIDOrName)
			if err != nil {
				return err
			}
			if server == nil {
				return fmt.Errorf("server not found: %s", serverIDOrName)
			}
			opts.Server = server
		}
		if automount {
			opts.Automount = &automount
		}
		if format != "" {
			opts.Format = &format
		}

		result, _, err := client.Volume().Create(ctx, opts)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}
		if err := waiter.WaitForActions(ctx, result.NextActions); err != nil {
			return err
		}
		fmt.Printf("Volume %d created\n", result.Volume.ID)

		return nil
	},
}
