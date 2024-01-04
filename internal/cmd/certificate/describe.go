package certificate

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "certificate",
	ShortDescription:     "Describe an certificate",
	JSONKeyGetByID:       "certificate",
	JSONKeyGetByName:     "certificates",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Certificate().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		cert, _, err := s.Client().Certificate().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return cert, hcloud.SchemaFromCertificate(cert), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		cert := resource.(*hcloud.Certificate)
		cmd.Printf("ID:\t\t\t%d\n", cert.ID)
		cmd.Printf("Name:\t\t\t%s\n", cert.Name)
		cmd.Printf("Type:\t\t\t%s\n", cert.Type)
		cmd.Printf("Fingerprint:\t\t%s\n", cert.Fingerprint)
		cmd.Printf("Created:\t\t%s (%s)\n", util.Datetime(cert.Created), humanize.Time(cert.Created))
		cmd.Printf("Not valid before:\t%s (%s)\n", util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
		cmd.Printf("Not valid after:\t%s (%s)\n", util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))
		if cert.Status != nil {
			cmd.Printf("Status:\n")
			cmd.Printf("  Issuance:\t%s\n", cert.Status.Issuance)
			cmd.Printf("  Renewal:\t%s\n", cert.Status.Renewal)
			if cert.Status.IsFailed() {
				cmd.Printf("  Failure reason: %s\n", cert.Status.Error.Message)
			}
		}
		cmd.Printf("Domain names:\n")
		for _, domainName := range cert.DomainNames {
			cmd.Printf("  - %s\n", domainName)
		}
		cmd.Print("Labels:\n")
		if len(cert.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range cert.Labels {
				cmd.Printf("  %s:\t%s\n", key, value)
			}
		}
		cmd.Println("Used By:")
		if len(cert.UsedBy) == 0 {
			cmd.Println("  Certificate unused")
		} else {
			for _, ub := range cert.UsedBy {
				cmd.Printf("  - Type: %s\n", ub.Type)
				// Currently certificates can be only attached to load balancers.
				// If we ever get something that is not a load balancer fall back
				// to printing the ID.
				if ub.Type != hcloud.CertificateUsedByRefTypeLoadBalancer {
					cmd.Printf("  - ID: %d\n", ub.ID)
					continue
				}
				cmd.Printf("  - Name: %s\n", s.Client().LoadBalancer().LoadBalancerName(ub.ID))
			}
		}
		return nil
	},
}
