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
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", cert.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", cert.Name)
		_, _ = fmt.Fprintf(out, "Type:\t%s\n", cert.Type)
		_, _ = fmt.Fprintf(out, "Fingerprint:\t%s\n", cert.Fingerprint)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(cert.Created), humanize.Time(cert.Created))
		_, _ = fmt.Fprintf(out, "Not valid before:\t%s (%s)\n", util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore))
		_, _ = fmt.Fprintf(out, "Not valid after:\t%s (%s)\n", util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))
		if cert.Status != nil {
			_, _ = fmt.Fprintf(out, "Status:\t\n")
			_, _ = fmt.Fprintf(out, "  Issuance:\t%s\n", cert.Status.Issuance)
			_, _ = fmt.Fprintf(out, "  Renewal:\t%s\n", cert.Status.Renewal)
			if cert.Status.IsFailed() {
				_, _ = fmt.Fprintf(out, "  Failure reason:\t%s\n", cert.Status.Error.Message)
			}
		}

		if len(cert.DomainNames) == 0 {
			_, _ = fmt.Fprintf(out, "Domain names:\tNo Domain names\n")
		} else {
			_, _ = fmt.Fprintf(out, "Domain names:\t\n")
			for _, domainName := range cert.DomainNames {
				_, _ = fmt.Fprintf(out, "\t- %s\n", domainName)
			}

		}

		util.DescribeLabels(out, cert.Labels, "")

		if len(cert.UsedBy) == 0 {
			_, _ = fmt.Fprintf(out, "Used By:\tCertificate unused\n")
		} else {
			_, _ = fmt.Fprintf(out, "Used By:\t\n")
			for _, ub := range cert.UsedBy {
				_, _ = fmt.Fprintf(out, "  - Type:\t%s\n", ub.Type)
				// Currently certificates can be only attached to load balancers.
				// If we ever get something that is not a load balancer fall back
				// to printing the ID.
				if ub.Type != hcloud.CertificateUsedByRefTypeLoadBalancer {
					_, _ = fmt.Fprintf(out, "  - ID:\t%d\n", ub.ID)
					continue
				}
				_, _ = fmt.Fprintf(out, "  - Name:\t%s\n", s.Client().LoadBalancer().LoadBalancerName(ub.ID))
			}
		}
		return nil
	},
}
