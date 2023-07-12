package certificate

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/cobra"

	"github.com/dustin/go-humanize"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "certificate",
	ShortDescription:     "Describe an certificate",
	JSONKeyGetByID:       "certificate",
	JSONKeyGetByName:     "certificates",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Certificate().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Certificate().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}) error {
		cert := resource.(*hcloud.Certificate)
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
				if ub.Type != hcloud.CertificateUsedByRefTypeLoadBalancer {
					fmt.Printf("  - ID: %d", ub.ID)
					continue
				}
				fmt.Printf("  - Name: %s", client.LoadBalancer().LoadBalancerName(ub.ID))
			}
		}
		return nil
	},
}
