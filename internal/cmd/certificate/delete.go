package certificate

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.Certificate]{
	ResourceNameSingular: "Certificate",
	ResourceNamePlural:   "Certificates",
	ShortDescription:     "Delete a Certificate",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Certificate().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Certificate, *hcloud.Response, error) {
		return s.Client().Certificate().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, certificate *hcloud.Certificate) ([]*hcloud.Action, error) {
		_, err := s.Client().Certificate().Delete(cmd.Context(), certificate)
		return nil, err
	},
}
