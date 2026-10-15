package server

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
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
		cmd.Flags().StringP("output", "o", "screenshot.png", "Output file for the screenshot ('-' for stdout; must be a .png file otherwise)")
		cmd.Flags().Duration("timeout", 10*time.Second, "Timeout for taking the screenshot (0s disables the timeout)")
		cmd.Flags().Bool("overwrite", false, "Overwrite the output file if it exists already (true, false)")
		return cmd
	},
	Experimental: base.ExperimentalWrapper("The server screenshot command", "experimental", ""),
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		output, _ := cmd.Flags().GetString("output")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		toStdout := output == "-"

		if toStdout && s.Terminal().StdoutIsTerminal() {
			return fmt.Errorf("interactive tty detected. use --output <file> or redirect stdout")
		}

		if !toStdout && !overwrite {
			if _, err := os.Stat(output); err == nil {
				return fmt.Errorf("output file %s already exists", output)
			}
		}

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

		if !toStdout && !strings.HasSuffix(strings.ToLower(output), ".png") {
			cmd.Printf("Warning: Only .png output file are supported\n")
		}

		var w io.Writer
		if toStdout {
			w = cmd.OutOrStdout()
		} else {
			f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return err
			}
			defer f.Close()
			w = f
		}

		isDebug, err := config.OptionDebug.Get(s.Config())
		if err != nil {
			return err
		}
		var debugOut io.Writer
		if isDebug {
			debugOut = cmd.ErrOrStderr()
		}

		if err := util.TakeScreenshot(s, result.WSSURL, w, debugOut, timeout); err != nil {
			return err
		}

		if !toStdout {
			cmd.Printf("Screenshot captured and saved to '%s'\n", output)
		}
		return nil
	},
}
