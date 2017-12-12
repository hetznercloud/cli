package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newImageDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "describe [flags] <id>",
		Short:            "Describe an Image",
		Args:             cobra.ExactArgs(1),
		TraverseChildren: true,
		RunE:             cli.wrap(runImageDescribe),
	}
	return cmd
}

func runImageDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Image ID")
	}

	ctx := context.Background()

	image, _, err := cli.Client().Image.Get(ctx, id)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("Image not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", image.ID)
	fmt.Printf("Type:\t\t%s\n", image.Type)
	fmt.Printf("Name:\t\t%s\n", image.Name)
	fmt.Printf("Description:\t%s\n", image.Description)
	fmt.Printf("Status:\t\t%s\n", image.Status)
	fmt.Printf("ImageSize:\t%.1f GB\n", image.ImageSize)
	fmt.Printf("DiskSize:\t%.0f GB\n", image.DiskSize)
	fmt.Printf("Created:\t%s\n", image.Created)
	fmt.Printf("Version:\t%s\n", image.Version)
	fmt.Printf("OSFlavor:\t%s\n", image.OSFlavor)
	fmt.Printf("OSVersion:\t%s\n", image.OSVersion)
	fmt.Printf("RapidDeploy:\t%t\n", image.RapidDeploy)

	return nil
}
