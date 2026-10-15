package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var VNCCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "vnc [options] <server>",
			Short: "Forward the Server VNC console to a local port",
			Long: `Forward the Server VNC console to a local port, so it can be accessed with any VNC client.

Note that the local VNC server has no password protection.`,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().String("address", "127.0.0.1", "Local address to listen on")
		cmd.Flags().IntP("port", "p", 5900, "Local port to listen on (0 picks a random free port)")
		return cmd
	},
	Experimental: base.ExperimentalWrapper("The server vnc command", "experimental", ""),
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")

		if port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %d", port)
		}

		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}

		// capture ctrl + c to stop forwarding
		ctx, stop := signal.NotifyContext(s, os.Interrupt)

		var lc net.ListenConfig
		listener, err := lc.Listen(ctx, "tcp", net.JoinHostPort(address, strconv.Itoa(port)))
		if err != nil {
			stop()
			return fmt.Errorf("could not listen on local address: %w", err)
		}

		var wg sync.WaitGroup
		defer func() {
			stop()
			wg.Wait()
		}()
		context.AfterFunc(ctx, func() {
			_ = listener.Close()
		})

		cmd.Printf("Forwarding VNC console of Server %d to %s, press Ctrl+C to stop\n", server.ID, listener.Addr())

		for {
			conn, err := listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					// ctrl + c pressed
					return nil
				}
				return err
			}
			wg.Go(func() {
				if err := forwardVNC(ctx, s, cmd, server, conn); err != nil && ctx.Err() == nil {
					cmd.PrintErrf("Error: connection from %s: %v\n", conn.RemoteAddr(), err)
				}
			})
		}
	},
}

func forwardVNC(ctx context.Context, s state.State, cmd *cobra.Command, server *hcloud.Server, conn net.Conn) error {
	defer conn.Close()

	result, _, err := s.Client().Server().RequestConsole(ctx, server)
	if err != nil {
		return err
	}
	if err := s.WaitForActions(ctx, cmd, result.Action); err != nil {
		return err
	}

	ws, err := util.DialWebsocket(ctx, result.WSSURL)
	if err != nil {
		return err
	}

	cmd.Printf("Accepted connection from %s\n", conn.RemoteAddr())
	defer cmd.Printf("Closed connection from %s\n", conn.RemoteAddr())

	return util.ProxyWebsocket(ctx, conn, ws)
}
