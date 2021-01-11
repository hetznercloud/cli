package volume

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newUpdateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] VOLUME",
		Short:                 "Update a volume",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.VolumeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runVolumeUpdate),
	}

	cmd.Flags().String("name", "", "Volume name")

	return cmd
}

func runVolumeUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	volume, _, err := cli.Client().Volume.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.VolumeUpdateOpts{
		Name: name,
	}
	if opts.Name == "" {
		return errors.New("no updates")
	}

	_, _, err = cli.Client().Volume.Update(cli.Context, volume, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Volume %d updated\n", volume.ID)
	return nil
}
