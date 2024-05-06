package server

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/ui"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ShutdownCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		const description = "Shuts down a Server gracefully by sending an ACPI shutdown request. " +
			"The Server operating system must support ACPI and react to the request, " +
			"otherwise the Server will not shut down. Use the --wait flag to wait for the " +
			"server to shut down before returning."

		cmd := &cobra.Command{
			Use:                   "shutdown [options] <server>",
			Short:                 "Shutdown a server",
			Long:                  description,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("wait", false, "Wait for the server to shut down before exiting")
		cmd.Flags().Duration("wait-timeout", 30*time.Second, "Timeout for waiting for off state after shutdown")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {

		wait, _ := cmd.Flags().GetBool("wait")
		timeout, _ := cmd.Flags().GetDuration("wait-timeout")

		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		action, _, err := s.Client().Server().Shutdown(s, server)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd, s, action); err != nil {
			return err
		}

		cmd.Printf("Sent shutdown signal to server %d\n", server.ID)

		if wait {
			start := time.Now()

			interval, _ := cmd.Flags().GetDuration("poll-interval")
			if interval < time.Second {
				interval = time.Second
			}

			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			progress := ui.NewProgress(os.Stderr, "Waiting for server to shut down")
			for server.Status != hcloud.ServerStatusOff {
				if now := <-ticker.C; now.Sub(start) >= timeout {
					progress.SetError()
					return errors.New("failed to shut down server")
				}

				server, _, err = s.Client().Server().GetByID(s, server.ID)
				if err != nil {
					progress.SetError()
					return err
				}
			}

			progress.SetSuccess()

			cmd.Printf("Server %d shut down\n", server.ID)
		}

		return nil
	},
}
