package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func newVolumeDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] VOLUME",
		Short:                 "Describe a volume",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runVolumeDescribe),
	}
	return cmd
}

func runVolumeDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	volume, _, err := cli.Client().Volume.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if volume == nil {
		return fmt.Errorf("volume not found: %s", args[0])
	}

	fmt.Printf("ID:\t\t%d\n", volume.ID)
	fmt.Printf("Name:\t\t%s\n", volume.Name)
	fmt.Printf("Size:\t\t%s\n", humanize.Bytes(uint64(volume.Size*humanize.GByte)))
	fmt.Printf("Linux Device:   %s\n", volume.LinuxDevice)
	fmt.Printf("Location:\n")
	fmt.Printf("  Name:\t\t%s\n", volume.Location.Name)
	fmt.Printf("  Description:\t%s\n", volume.Location.Description)
	fmt.Printf("  Country:\t%s\n", volume.Location.Country)
	fmt.Printf("  City:\t\t%s\n", volume.Location.City)
	fmt.Printf("  Latitude:\t%f\n", volume.Location.Latitude)
	fmt.Printf("  Longitude:\t%f\n", volume.Location.Longitude)
	if volume.Server != nil {
		server, _, err := cli.Client().Server.GetByID(cli.Context, volume.Server.ID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %d", volume.Server.ID)
		}
		fmt.Printf("Server:\n")
		fmt.Printf("  ID:\t\t%d\n", server.ID)
		fmt.Printf("  Name:\t\t%s\n", server.Name)
	} else {
		fmt.Print("Server:\n  Not attached\n")
	}
	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", yesno(volume.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(volume.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range volume.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}
