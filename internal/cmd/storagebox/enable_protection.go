package storagebox

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.StorageBoxChangeProtectionOpts, error) {
	opts := hcloud.StorageBoxChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = &enable
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
	storageBox *hcloud.StorageBox, enable bool, opts hcloud.StorageBoxChangeProtectionOpts) error {

	if opts.Delete == nil {
		return nil
	}

	action, _, err := s.Client().StorageBox().ChangeProtection(s, storageBox, opts)
	if err != nil {
		return err
	}

	if err := s.WaitForActions(s, cmd, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for Storage Box %d\n", storageBox.ID)
	} else {
		cmd.Printf("Resource protection disabled for Storage Box %d\n", storageBox.ID)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection <storage-box> delete",
			Args:  util.ValidateLenient,
			Short: "Enable resource protection for a Storage Box",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(true, args[1:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, storageBox, true, opts)
	},
	Experimental: experimental.StorageBoxes,
}
