package certificate

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "certificate",
	ShortDescriptionAdd:    "Add a label to an certificate",
	ShortDescriptionRemove: "Remove a label from an certificate",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Certificate().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Certificate().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		certificate, _, err := s.Client().Certificate().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if certificate == nil {
			return nil, fmt.Errorf("certificate not found: %s", idOrName)
		}
		return certificate, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		cert := resource.(*hcloud.Certificate)
		opts := hcloud.CertificateUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Certificate().Update(s, cert, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		cert := resource.(*hcloud.Certificate)
		return cert.Labels
	},
	GetIDOrName: func(resource any) string {
		cert := resource.(*hcloud.Certificate)
		return strconv.FormatInt(cert.ID, 10)
	},
}
