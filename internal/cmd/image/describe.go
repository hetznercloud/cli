package image

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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
			fmt.Fprintln(os.Stderr, "INFO: This command only returns x86 images by default. Explicitly set the --architecture=x86|arm flag to hide this message.")
		}

		img, _, err := s.Client().Image().GetForArchitecture(s, idOrName, hcloud.Architecture(arch))
		if err != nil {
			return nil, nil, err
		}
		return img, hcloud.SchemaFromImage(img), nil
	},
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, image *hcloud.Image) error {
		fmt.Fprint(out, DescribeImage(image))
		return nil
	},
}

func DescribeImage(image *hcloud.Image) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ID:\t%d\n", image.ID)
	fmt.Fprintf(&sb, "Type:\t%s\n", image.Type)
	fmt.Fprintf(&sb, "Status:\t%s\n", image.Status)
	fmt.Fprintf(&sb, "Name:\t%s\n", util.NA(image.Name))
	fmt.Fprintf(&sb, "Created:\t%s (%s)\n", util.Datetime(image.Created), humanize.Time(image.Created))
	if !image.Deprecated.IsZero() {
		fmt.Fprintf(&sb, "Deprecated:\t%s (%s)\n", util.Datetime(image.Deprecated), humanize.Time(image.Deprecated))
	}
	fmt.Fprintf(&sb, "Description:\t%s\n", image.Description)
	if image.ImageSize != 0 {
		fmt.Fprintf(&sb, "Image size:\t%.2f GB\n", image.ImageSize)
	} else {
		fmt.Fprintf(&sb, "Image size:\t%s\n", util.NA(""))
	}
	fmt.Fprintf(&sb, "Disk size:\t%.0f GB\n", image.DiskSize)
	fmt.Fprintf(&sb, "OS flavor:\t%s\n", image.OSFlavor)
	fmt.Fprintf(&sb, "OS version:\t%s\n", util.NA(image.OSVersion))
	fmt.Fprintf(&sb, "Architecture:\t%s\n", image.Architecture)
	fmt.Fprintf(&sb, "Rapid deploy:\t%s\n", util.YesNo(image.RapidDeploy))
	fmt.Fprintf(&sb, "Protection:\t\n")
	fmt.Fprintf(&sb, "  Delete:\t%s\n", util.YesNo(image.Protection.Delete))

	util.DescribeLabels(&sb, image.Labels, "")

	if !image.Deprecated.IsZero() {
		fmt.Fprintf(&sb, "\nAttention: This Image is deprecated and will be removed in the future.\n")
	}
	return sb.String()
}
