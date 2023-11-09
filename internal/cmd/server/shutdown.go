package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ShutdownCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {

		const description = "Shuts down a Server gracefully by sending an ACPI shutdown request. " +
			"The Server operating system must support ACPI and react to the request, " +
			"otherwise the Server will not shut down. Use the --wait flag to wait for the " +
			"server to shut down before returning."

		cmd := &cobra.Command{
			Use:                   "shutdown [FLAGS] SERVER",
			Short:                 "Shutdown a server",
			Long:                  description,
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("wait", false, "Wait for the server to shut down before exiting")
		cmd.Flags().Duration("wait-timeout", 30*time.Second, "Timeout for waiting for off state after shutdown")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {

		wait, _ := cmd.Flags().GetBool("wait")
		timeout, _ := cmd.Flags().GetDuration("wait-timeout")

		idOrName := args[0]
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		action, _, err := client.Server().Shutdown(ctx, server)
		if err != nil {
			return err
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}

		cmd.Printf("Sent shutdown signal to server %d\n", server.ID)

		if wait {
			start := time.Now()
			errCh := make(chan error)

			interval, _ := cmd.Flags().GetDuration("poll-interval")
			if interval < time.Second {
				interval = time.Second
			}

			go func() {
				defer close(errCh)

				ticker := time.NewTicker(interval)
				defer ticker.Stop()

				for server.Status != hcloud.ServerStatusOff {
					if now := <-ticker.C; now.Sub(start) >= timeout {
						errCh <- errors.New("failed to shut down server")
						return
					}
					server, _, err = client.Server().GetByID(ctx, server.ID)
					if err != nil {
						errCh <- err
						return
					}
				}

				errCh <- nil
			}()

			if err := state.DisplayProgressCircle(errCh, "Waiting for server to shut down"); err != nil {
				return err
			}

			cmd.Printf("Server %d shut down\n", server.ID)
		}

		return nil
	},
}
