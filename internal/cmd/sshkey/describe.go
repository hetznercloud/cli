package sshkey

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] SSHKEY",
		Short:                 "Describe a SSH key",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.SSHKeyNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runSSHKeyDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runSSHKeyDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

	sshKey, resp, err := cli.Client().SSHKey.Get(cli.Context, args[0])
	if err != nil {
		return err
	}
	if sshKey == nil {
		return fmt.Errorf("SSH key not found: %s", args[0])
	}

	switch {
	case outputFlags.IsSet("json"):
		return sshKeyDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(sshKey, outputFlags["format"][0])
	default:
		return sshKeyDescribeText(cli, sshKey)
	}
}

func sshKeyDescribeText(cli *state.State, sshKey *hcloud.SSHKey) error {
	fmt.Printf("ID:\t\t%d\n", sshKey.ID)
	fmt.Printf("Name:\t\t%s\n", sshKey.Name)
	fmt.Printf("Created:\t%s (%s)\n", util.Datetime(sshKey.Created), humanize.Time(sshKey.Created))
	fmt.Printf("Fingerprint:\t%s\n", sshKey.Fingerprint)
	fmt.Printf("Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))
	fmt.Print("Labels:\n")
	if len(sshKey.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range sshKey.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	return nil
}

func sshKeyDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if sshKey, ok := data["ssh_key"]; ok {
		return util.DescribeJSON(sshKey)
	}
	if sshKeys, ok := data["ssh_keys"].([]interface{}); ok {
		return util.DescribeJSON(sshKeys[0])
	}
	return util.DescribeJSON(data)
}
