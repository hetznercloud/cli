package certificate

import (
	"fmt"

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
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		certificate, _, err := s.Certificate().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if certificate == nil {
			return nil, 0, fmt.Errorf("certificate not found: %s", idOrName)
		}
		return certificate.Labels, certificate.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.CertificateUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Certificate().Update(s, &hcloud.Certificate{ID: id}, opts)
		return err
	},
}
