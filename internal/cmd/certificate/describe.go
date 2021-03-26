package certificate

import (
	"encoding/json"
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] CERTIFICATE",
		Short:                 "Describe a certificate",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.CertificateNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDescribe),
	}
	output.AddFlag(cmd, output.OptionJSON(), output.OptionFormat())
	return cmd
}

func runDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := output.FlagsForCommand(cmd)

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
		return describeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(cert, outputFlags["format"][0])
	default:
		return describeText(cli, cert)
	}
}

func describeJSON(resp *hcloud.Response) error {
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

func describeText(cli *state.State, cert *hcloud.Certificate) error {
	fmt.Printf("ID:\t\t\t%d\n", cert.ID)
	fmt.Printf("Name:\t\t\t%s\n", cert.Name)
	fmt.Printf("Type:\t\t\t%s\n", cert.Type)
	fmt.Printf("Fingerprint:\t\t%s\n", cert.Fingerprint)
	fmt.Printf("Created:\t\t%s (%s)\n", util.Datetime(cert.Created), humanize.Time(cert.Created))
	fmt.Printf("Not valid before:\t%s (%s)\n", util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
	fmt.Printf("Not valid after:\t%s (%s)\n", util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))
	if cert.Status != nil {
		fmt.Printf("Status:\n")
		fmt.Printf("  Issuance: %s\n", cert.Status.Issuance)
		fmt.Printf("  Renewal: %s\n", cert.Status.Renewal)
		if cert.Status.IsFailed() {
			fmt.Printf("  Failure reason: %s\n", cert.Status.Error.Message)
		}
	}
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
	fmt.Println("Used By:")
	if len(cert.UsedBy) == 0 {
		fmt.Println("  Certificate unused")
	} else {
		for _, ub := range cert.UsedBy {
			fmt.Printf("  - Type: %s", ub.Type)
			// Currently certificates can be only attached to load balancers.
			// If we ever get something that is not a load balancer fall back
			// to printing the ID.
			if ub.Type != "load_balancer" {
				fmt.Printf("  - ID: %d", ub.ID)
				continue
			}
			fmt.Printf("  - Name: %s", cli.LoadBalancerName(ub.ID))
		}
	}
	return nil
}
