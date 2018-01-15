package cli

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func newImageDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] IMAGE",
		Short:                 "Describe an image",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runImageDescribe),
	}
	return cmd
}

func runImageDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	image, _, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}

	fmt.Printf("ID:\t\t%d\n", image.ID)
	fmt.Printf("Type:\t\t%s\n", image.Type)
	fmt.Printf("Status:\t\t%s\n", image.Status)
	fmt.Printf("Name:\t\t%s\n", na(image.Name))
	fmt.Printf("Description:\t%s\n", image.Description)
	if image.ImageSize != 0 {
		fmt.Printf("Image size:\t%.1f GB\n", image.ImageSize)
	} else {
		fmt.Printf("Image size:\t%s\n", na(""))
	}
	fmt.Printf("Disk size:\t%.0f GB\n", image.DiskSize)
	fmt.Printf("Created:\t%s (%s)\n", image.Created, humanize.Time(image.Created))
	fmt.Printf("OS flavor:\t%s\n", image.OSFlavor)
	fmt.Printf("OS version:\t%s\n", na(image.OSVersion))
	fmt.Printf("Rapid deploy:\t%s\n", yesno(image.RapidDeploy))

	return nil
}
