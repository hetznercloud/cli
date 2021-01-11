package cmds

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newCertificateDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] CERTIFICATE",
		Short:                 "Describe a certificate",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.CertificateNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runCertificateDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runCertificateDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	cert, resp, err := cli.Client().Certificate.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if cert == nil {
		return fmt.Errorf("certificate not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return certificateDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(cert, outputFlags["format"][0])
	default:
		return certificateDescribeText(cli, cert)
	}
}

func certificateDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if server, ok := data["certificate"]; ok {
		return util.DescribeJSON(server)
	}
	if servers, ok := data["certificates"].([]interface{}); ok {
		return util.DescribeJSON(servers[0])
	}
	return util.DescribeJSON(data)
}

func certificateDescribeText(cli *state.State, cert *hcloud.Certificate) error {
	fmt.Printf("ID:\t\t\t%d\n", cert.ID)
	fmt.Printf("Name:\t\t\t%s\n", cert.Name)
	fmt.Printf("Fingerprint:\t\t%s\n", cert.Fingerprint)
	fmt.Printf("Created:\t\t%s (%s)\n", util.Datetime(cert.Created), humanize.Time(cert.Created))
	fmt.Printf("Not valid before:\t%s (%s)\n", util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
	fmt.Printf("Not valid after:\t%s (%s)\n", util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))
	fmt.Printf("Domain names:\n")
	for _, domainName := range cert.DomainNames {
		fmt.Printf("  - %s\n", domainName)
	}
	fmt.Print("Labels:\n")
	if len(cert.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range cert.Labels {
			fmt.Printf("  %s:\t%s\n", key, value)
		}
	}
	return nil
}
