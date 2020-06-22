package cli

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newCertificateDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] CERTIFICATE",
		Short:                 "Describe a certificate",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runCertificateDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runCertificateDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
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
		return describeFormat(cert, outputFlags["format"][0])
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
		return describeJSON(server)
	}
	if servers, ok := data["certificates"].([]interface{}); ok {
		return describeJSON(servers[0])
	}
	return describeJSON(data)
}

func certificateDescribeText(cli *CLI, cert *hcloud.Certificate) error {
	fmt.Printf("ID:\t\t\t%d\n", cert.ID)
	fmt.Printf("Name:\t\t\t%s\n", cert.Name)
	fmt.Printf("Fingerprint:\t\t%s\n", cert.Fingerprint)
	fmt.Printf("Created:\t\t%s (%s)\n", datetime(cert.Created), humanize.Time(cert.Created))
	fmt.Printf("Not valid before:\t%s (%s)\n", datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
	fmt.Printf("Not valid after:\t%s (%s)\n", datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))
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
