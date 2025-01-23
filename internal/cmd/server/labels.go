package server

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	state "github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "server",
	ShortDescriptionAdd:    "Add a label to a server",
	ShortDescriptionRemove: "Remove a label from a server",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Server().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Server().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if server == nil {
			return nil, fmt.Errorf("server not found: %s", idOrName)
		}
		return server, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		server := resource.(*hcloud.Server)
		opts := hcloud.ServerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Server().Update(s, server, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		server := resource.(*hcloud.Server)
		return server.Labels
	},
	GetIDOrName: func(resource any) string {
		server := resource.(*hcloud.Server)
		return strconv.FormatInt(server.ID, 10)
	},
}
