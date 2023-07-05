package certificate

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var labelCmds = base.LabelCmds{
	ResourceNameSingular:   "certificate",
	ShortDescriptionAdd:    "Add a label to an certificate",
	ShortDescriptionRemove: "Remove a label from an certificate",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Certificate().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Certificate().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error) {
		certificate, _, err := client.Certificate().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if certificate == nil {
			return nil, 0, fmt.Errorf("certificate not found: %s", idOrName)
		}
		return certificate.Labels, certificate.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error {
		opts := hcloud.CertificateUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.Certificate().Update(ctx, &hcloud.Certificate{ID: id}, opts)
		return err
	},
}
