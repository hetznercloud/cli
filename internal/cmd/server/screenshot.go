package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
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
		cmd.Flags().Duration("watch", 0, "Take a screenshot at the given interval (minimum 1s)")

		cmd.MarkFlagsMutuallyExclusive("timeout", "watch")
		return cmd
	},
	Experimental: base.ExperimentalWrapper("The server screenshot command", "experimental", ""),
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		output, _ := cmd.Flags().GetString("output")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		watch, _ := cmd.Flags().GetDuration("watch")
		isWatch := cmd.Flags().Changed("watch")
		toStdout := output == "-"

		if isWatch {
			if watch < time.Second {
				return fmt.Errorf("--watch interval must be at least 1s")
			}
			if toStdout {
				return fmt.Errorf("--watch cannot be used with output to stdout")
			}
		}

		if toStdout && s.Terminal().StdoutIsTerminal() {
			return fmt.Errorf("interactive tty detected. use --output <file> or redirect stdout")
		}

		if !toStdout && !overwrite {
			firstOutput := output
			if isWatch {
				firstOutput = util.IndexedFilePath(output, 1)
			}
			if _, err := os.Stat(firstOutput); err == nil {
				return fmt.Errorf("output file %s already exists", firstOutput)
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

		isDebug, err := config.OptionDebug.Get(s.Config())
		if err != nil {
			return err
		}
		var debugOut io.Writer
		if isDebug {
			debugOut = cmd.ErrOrStderr()
		}

		if isWatch {
			return watchScreenshots(s, cmd, result.WSSURL, output, overwrite, watch, debugOut)
		}

		ctx := context.Context(s)
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		session, err := util.OpenScreenshotSession(ctx, result.WSSURL, debugOut)
		if err != nil {
			return err
		}
		defer session.Close()

		var buf bytes.Buffer
		if err := session.Capture(ctx, &buf); err != nil {
			return err
		}

		if toStdout {
			_, err := cmd.OutOrStdout().Write(buf.Bytes())
			return err
		}
		if err := writeScreenshotFile(output, buf.Bytes(), overwrite); err != nil {
			return err
		}
		cmd.Printf("Screenshot captured and saved to '%s'\n", output)
		return nil
	},
}

func watchScreenshots(
	ctx context.Context, cmd *cobra.Command, wsURL, output string, overwrite bool,
	interval time.Duration, debugOut io.Writer,
) error {
	// capture ctrl + c to stop taking screenshots
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	session, err := util.OpenScreenshotSession(ctx, wsURL, debugOut)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return err
	}
	defer session.Close()

	cmd.Printf("Taking a screenshot every %s, press Ctrl+C to stop\n", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for idx := 1; ; idx++ {
		var buf bytes.Buffer
		if err := session.Capture(ctx, &buf); err != nil {
			if ctx.Err() != nil {
				// ctrl + c signal received
				return nil
			}
			return err
		}

		path := util.IndexedFilePath(output, idx)
		if err := writeScreenshotFile(path, buf.Bytes(), overwrite); err != nil {
			return err
		}
		cmd.Printf("Screenshot captured and saved to '%s'\n", path)

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func writeScreenshotFile(path string, data []byte, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("output file %s already exists", path)
		}
	}
	return os.WriteFile(path, data, 0600)
}
