package certificate

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "certificate",
	ShortDescription:     "Update a certificate",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.Certificate().Get(ctx, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Certificate Name")
	},
	Update: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		certificate := resource.(*hcloud.Certificate)
		updOpts := hcloud.CertificateUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := client.Certificate().Update(ctx, certificate, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
