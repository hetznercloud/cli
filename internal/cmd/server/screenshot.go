package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var ScreenshotCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "screenshot [options] <server>",
			Short:                 "Take a screenshot of the Server VNC console",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().StringP("output", "o", "screenshot.png", "Output file for the screenshot (must be a .png file).")
		cmd.Flags().Duration("timeout", 10*time.Second, "Timeout for taking the screenshot (0s disables the timeout)")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		output, _ := cmd.Flags().GetString("output")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		result, _, err := s.Client().Server().RequestConsole(s, server)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return err
		}

		if !strings.HasSuffix(strings.ToLower(output), ".png") {
			cmd.Printf("Warning: Only .png output file are supported\n")
		}

		return util.TakeScreenshot(s, result.WSSURL, output, timeout)
	},
}
