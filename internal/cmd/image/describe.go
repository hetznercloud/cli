package image

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] IMAGE",
		Short:                 "Describe an image",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ImageNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	idOrName := args[0]
	image, resp, err := cli.Client().Image.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if image == nil {
		return fmt.Errorf("image not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return describeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(image, outputFlags["format"][0])
	default:
		return describeText(cli, image)
	}
}

func describeText(cli *state.State, image *hcloud.Image) error {
	fmt.Printf("ID:\t\t%d\n", image.ID)
	fmt.Printf("Type:\t\t%s\n", image.Type)
	fmt.Printf("Status:\t\t%s\n", image.Status)
	fmt.Printf("Name:\t\t%s\n", util.NA(image.Name))
	fmt.Printf("Created:\t%s (%s)\n", util.Datetime(image.Created), humanize.Time(image.Created))
	if !image.Deprecated.IsZero() {
		fmt.Printf("Deprecated:\t%s (%s)\n", util.Datetime(image.Deprecated), humanize.Time(image.Deprecated))
	}
	fmt.Printf("Description:\t%s\n", image.Description)
	if image.ImageSize != 0 {
		fmt.Printf("Image size:\t%.2f GB\n", image.ImageSize)
	} else {
		fmt.Printf("Image size:\t%s\n", util.NA(""))
	}
	fmt.Printf("Disk size:\t%.0f GB\n", image.DiskSize)
	fmt.Printf("OS flavor:\t%s\n", image.OSFlavor)
	fmt.Printf("OS version:\t%s\n", util.NA(image.OSVersion))
	fmt.Printf("Rapid deploy:\t%s\n", util.YesNo(image.RapidDeploy))
	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", util.YesNo(image.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(image.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range image.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	if !image.Deprecated.IsZero() {
		fmt.Printf("\nAttention: This image is deprecated and will be removed in the future.\n")
	}
	return nil
}

func describeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if image, ok := data["image"]; ok {
		return util.DescribeJSON(image)
	}
	if images, ok := data["images"].([]interface{}); ok {
		return util.DescribeJSON(images[0])
	}
	return util.DescribeJSON(data)
}
