package server

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var SSHCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "ssh [options] <server> [command...]",
			Short:                 "Spawn an SSH connection for the server",
			Args:                  util.ValidateLenient,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().SetInterspersed(false) // To make "hcloud server ssh <server> uname -a" execute "uname -a"
		cmd.Flags().Bool("ipv6", false, "Establish SSH connection to IPv6 address")
		cmd.Flags().StringP("user", "u", "root", "Username for SSH connection")
		cmd.Flags().IntP("port", "p", 22, "Port for SSH connection")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", idOrName)
		}

		useIPv6, _ := cmd.Flags().GetBool("ipv6")
		user, _ := cmd.Flags().GetString("user")
		port, _ := cmd.Flags().GetInt("port")

		ipAddress := server.PublicNet.IPv4.IP
		if server.PublicNet.IPv4.IsUnspecified() || useIPv6 {
			if server.PublicNet.IPv6.IsUnspecified() {
				return fmt.Errorf("server %s does not have a assigned primary ipv4 or ipv6", idOrName)
			}
			ipAddress = server.PublicNet.IPv6.Network.IP
			// increment last byte to get the ::1 IP, which is routed
			ipAddress[15]++
		}

		sshArgs := []string{"-l", user, "-p", strconv.Itoa(port), ipAddress.String()}
		sshCommand := exec.Command(s.Config().SSHPath(), append(sshArgs, args[1:]...)...)
		sshCommand.Stdin = os.Stdin
		sshCommand.Stdout = os.Stdout
		sshCommand.Stderr = os.Stderr

		if err := sshCommand.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				os.Exit(waitStatus.ExitStatus())
			} else {
				return err
			}
		}

		return nil
	},
}
