package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

func newServerSSHCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ssh [FLAGS] SERVER",
		Short:                 "Spawn an SSH connection for the server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerSSH),
	}
	cmd.Flags().Bool("ipv6", false, "Establish SSH connection to IPv6 address")
	cmd.Flags().String("user", "root", "Username for SSH connection")
	return cmd
}

func runServerSSH(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	useIPv6, _ := cmd.Flags().GetBool("ipv6")
	user, _ := cmd.Flags().GetString("user")

	ipAddress := server.PublicNet.IPv4.IP
	if useIPv6 {
		ipAddress = server.PublicNet.IPv6.Network.IP
		// increment last byte to get the ::1 IP, which is routed
		ipAddress[15]++
	}

	sshCommand := exec.Command("ssh", "-l", user, ipAddress.String())
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
}
