package server

import (
	"fmt"

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
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if server == nil {
			return nil, 0, fmt.Errorf("server not found: %s", idOrName)
		}
		return server.Labels, server.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.ServerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Server().Update(s, &hcloud.Server{ID: id}, opts)
		return err
	},
}
