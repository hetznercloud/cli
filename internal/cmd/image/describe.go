package image

import (
	"context"
	"fmt"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "image",
	ShortDescription:     "Describe an image",
	JSONKeyGetByID:       "image",
	JSONKeyGetByName:     "images",
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringP("architecture", "a", string(hcloud.ArchitectureX86), "architecture of the image, default is x86")
		cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))

	},
	NameSuggestions: func(c hcapi2.Client) func() []string { return c.Image().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		arch, err := cmd.Flags().GetString("architecture")
		if err != nil {
			return nil, nil, err
		}
		if !cmd.Flags().Changed("architecture") {
			_, _ = fmt.Fprintln(os.Stderr, "INFO: This command only returns x86 images by default. Explicitly set the --architecture=x86|arm flag to hide this message.")
		}
		return client.Image().GetForArchitecture(ctx, idOrName, hcloud.Architecture(arch))
	},
	PrintText: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, resource interface{}) error {
		image := resource.(*hcloud.Image)

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
		fmt.Printf("Architecture:\t%s\n", image.Architecture)
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
	},
}
