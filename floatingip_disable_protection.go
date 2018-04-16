package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPDisableProtectionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disable-protection [FLAGS] FLOATINGIP [PROTECTION LEVELS...]",
		Short:                 "Disable resource protection for a Floating IP",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPDisableProtection),
	}
	return cmd
}

func runFloatingIPDisableProtection(cli *CLI, cmd *cobra.Command, args []string) error {
	floatingIPID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}
	floatingIP := &hcloud.FloatingIP{ID: floatingIPID}

	var unknown []string
	opts := hcloud.FloatingIPChangeProtectionOpts{}
	for _, arg := range args[1:] {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = hcloud.Bool(false)
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	action, _, err := cli.Client().FloatingIP.ChangeProtection(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Resource protection disabled for Floating IP %d\n", floatingIP.ID)
	return nil
}
