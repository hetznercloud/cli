package certificate

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.Certificate]{
	ResourceNameSingular: "Certificate",
	ShortDescription:     "Describe a Certificate",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Certificate().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Certificate, any, error) {
		cert, _, err := s.Client().Certificate().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return cert, hcloud.SchemaFromCertificate(cert), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, cert *hcloud.Certificate) error {
		fmt.Fprintf(out, "ID:\t%d\n", cert.ID)
		fmt.Fprintf(out, "Name:\t%s\n", cert.Name)
		fmt.Fprintf(out, "Type:\t%s\n", cert.Type)
		fmt.Fprintf(out, "Fingerprint:\t%s\n", cert.Fingerprint)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(cert.Created), humanize.Time(cert.Created))
		fmt.Fprintf(out, "Not valid before:\t%s (%s)\n", util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
		fmt.Fprintf(out, "Not valid after:\t%s (%s)\n", util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))

		if cert.Status != nil {
			fmt.Fprintln(out)
			fmt.Fprintf(out, "Status:\n")
			fmt.Fprintf(out, "  Issuance:\t%s\n", cert.Status.Issuance)
			fmt.Fprintf(out, "  Renewal:\t%s\n", cert.Status.Renewal)
			if cert.Status.IsFailed() {
				fmt.Fprintf(out, "  Failure reason:\t%s\n", cert.Status.Error.Message)
			}
		}

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Domain names:\n")
		if len(cert.DomainNames) == 0 {
			fmt.Fprintf(out, "  No Domain names\n")
		} else {
			for _, domainName := range cert.DomainNames {
				fmt.Fprintf(out, "  - %s\n", domainName)
			}

		}

		fmt.Fprintln(out)
		util.DescribeLabels(out, cert.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Used By:\n")
		if len(cert.UsedBy) == 0 {
			fmt.Fprintf(out, "  Certificate unused\n")
		} else {
			for _, ub := range cert.UsedBy {
				fmt.Fprintf(out, "  - Type:\t%s\n", ub.Type)
				// Currently certificates can be only attached to load balancers.
				// If we ever get something that is not a load balancer fall back
				// to printing the ID.
				if ub.Type != hcloud.CertificateUsedByRefTypeLoadBalancer {
					fmt.Fprintf(out, "  - ID:\t%d\n", ub.ID)
					continue
				}
				fmt.Fprintf(out, "  - Name:\t%s\n", s.Client().LoadBalancer().LoadBalancerName(ub.ID))
			}
		}
		return nil
	},
}
