package certificate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.Certificate]{
	ResourceNameSingular:   "Certificate",
	ShortDescriptionAdd:    "Add a label to a Certificate",
	ShortDescriptionRemove: "Remove a label from a Certificate",
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Certificate().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.Certificate().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.Certificate, error) {
		certificate, _, err := s.Client().Certificate().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if certificate == nil {
			return nil, fmt.Errorf("Certificate not found: %s", idOrName)
		}
		return certificate, nil
	},
	SetLabels: func(ctx context.Context, s state.State, cert *hcloud.Certificate, labels map[string]string) error {
		opts := hcloud.CertificateUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Certificate().Update(ctx, cert, opts)
		return err
	},
	GetLabels: func(cert *hcloud.Certificate) map[string]string {
		return cert.Labels
	},
	GetIDOrName: func(cert *hcloud.Certificate) string {
		return strconv.FormatInt(cert.ID, 10)
	},
}
