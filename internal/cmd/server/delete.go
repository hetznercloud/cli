package server

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Server",
	ResourceNamePlural:   "Servers",
	ShortDescription:     "Delete a server",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Server().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().Server().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) (*hcloud.Action, error) {
		server := resource.(*hcloud.Server)
		result, _, err := s.Client().Server().DeleteWithResult(s, server)
		if err != nil {
			return nil, err
		}
		return result.Action, nil
	},
}
