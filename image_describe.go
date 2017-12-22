package cli

import (
	"errors"
	"fmt"
	"strconv"

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
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid image id")
	}

	image, _, err := cli.Client().Image.GetByID(cli.Context, id)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", image.ID)
	fmt.Printf("Type:\t\t%s\n", image.Type)
	fmt.Printf("Status:\t\t%s\n", image.Status)
	fmt.Printf("Name:\t\t%s\n", na(image.Name))
	fmt.Printf("Description:\t%s\n", image.Description)
	fmt.Printf("Image size:\t%.1f GB\n", image.ImageSize)
	fmt.Printf("Disk size:\t%.0f GB\n", image.DiskSize)
	fmt.Printf("Created:\t%s (%s)\n", image.Created, humanize.Time(image.Created))
	fmt.Printf("OS flavor:\t%s\n", image.OSFlavor)
	fmt.Printf("OS version:\t%s\n", na(image.OSVersion))
	fmt.Printf("Rapid deploy:\t%s\n", yesno(image.RapidDeploy))

	return nil
}
