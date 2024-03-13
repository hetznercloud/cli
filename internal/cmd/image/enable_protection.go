package image

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.ImageChangeProtectionOpts, error) {

	opts := hcloud.ImageChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Ptr(enable)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(s state.State, cmd *cobra.Command,
	image *hcloud.Image, enable bool, opts hcloud.ImageChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := s.Client().Image().ChangeProtection(s, image, opts)
	if err != nil {
		return err
	}

	if err := s.ActionProgress(cmd, s, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for image %d\n", image.ID)
	} else {
		cmd.Printf("Resource protection disabled for image %d\n", image.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		return &cobra.Command{
			Use:   "enable-protection <image> <protection-level>...",
			Short: "Enable resource protection for an image",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Image().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		imageID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.New("invalid image ID")
		}
		image := &hcloud.Image{ID: imageID}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, image, true, opts)
	},
}
