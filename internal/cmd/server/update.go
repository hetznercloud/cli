package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd[*hcloud.Server]{
	ResourceNameSingular: "Server",
	ShortDescription:     "Update a Server",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Server().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.Server, *hcloud.Response, error) {
		return s.Client().Server().Get(cmd.Context(), idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "Server name")
	},
	Update: func(s state.State, cmd *cobra.Command, server *hcloud.Server, flags map[string]pflag.Value) error {
		updOpts := hcloud.ServerUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.Client().Server().Update(cmd.Context(), server, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}
