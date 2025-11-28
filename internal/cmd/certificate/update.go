package certificate

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.Certificate]{
	ResourceNameSingular: "Certificate",
	ShortDescription:     "Update a Certificate",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.Certificate, *hcloud.Response, error) {
		return s.Client().Certificate().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Certificate Name")
	},
	Update: func(s state.State, _ *cobra.Command, certificate *hcloud.Certificate, flags map[string]pflag.Value) error {
		updOpts := hcloud.CertificateUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Certificate().Update(s, certificate, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
