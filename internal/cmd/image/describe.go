package image

import (
	"fmt"
	"os"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Image]{
	ResourceNameSingular: "Image",
	ShortDescription:     "Describe an Image",
	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().StringP("architecture", "a", string(hcloud.ArchitectureX86), "architecture of the Image, default is x86")
		_ = cmd.RegisterFlagCompletionFunc("architecture", cmpl.SuggestCandidates(string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)))
	},
	NameSuggestions: func(c hcapi2.Client) func() []string { return c.Image().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Image, any, error) {
		_, err := strconv.ParseInt(idOrName, 10, 64)
		isID := err == nil

		arch, err := cmd.Flags().GetString("architecture")
		if err != nil {
			return nil, nil, err
		}

		if !isID && !cmd.Flags().Changed("architecture") {
			_, _ = fmt.Fprintln(os.Stderr, "INFO: This command only returns x86 images by default. Explicitly set the --architecture=x86|arm flag to hide this message.")
		}

		img, _, err := s.Client().Image().GetForArchitecture(s, idOrName, hcloud.Architecture(arch))
		if err != nil {
			return nil, nil, err
		}
		return img, hcloud.SchemaFromImage(img), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, image *hcloud.Image, _ base.DescribeWriter) error {
		cmd.Printf("ID:\t\t%d\n", image.ID)
		cmd.Printf("Type:\t\t%s\n", image.Type)
		cmd.Printf("Status:\t\t%s\n", image.Status)
		cmd.Printf("Name:\t\t%s\n", util.NA(image.Name))
		cmd.Printf("Created:\t%s (%s)\n", util.Datetime(image.Created), humanize.Time(image.Created))
		if !image.Deprecated.IsZero() {
			cmd.Printf("Deprecated:\t%s (%s)\n", util.Datetime(image.Deprecated), humanize.Time(image.Deprecated))
		}
		cmd.Printf("Description:\t%s\n", image.Description)
		if image.ImageSize != 0 {
			cmd.Printf("Image size:\t%.2f GB\n", image.ImageSize)
		} else {
			cmd.Printf("Image size:\t%s\n", util.NA(""))
		}
		cmd.Printf("Disk size:\t%.0f GB\n", image.DiskSize)
		cmd.Printf("OS flavor:\t%s\n", image.OSFlavor)
		cmd.Printf("OS version:\t%s\n", util.NA(image.OSVersion))
		cmd.Printf("Architecture:\t%s\n", image.Architecture)
		cmd.Printf("Rapid deploy:\t%s\n", util.YesNo(image.RapidDeploy))
		cmd.Printf("Protection:\n")
		cmd.Printf("  Delete:\t%s\n", util.YesNo(image.Protection.Delete))

		cmd.Print("Labels:\n")
		if len(image.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(image.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		if !image.Deprecated.IsZero() {
			cmd.Printf("\nAttention: This Image is deprecated and will be removed in the future.\n")
		}
		return nil
	},
}
