package server

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	state "github.com/hetznercloud/cli/internal/state"
)

var LabelCmds = base.LabelCmds[*hcloud.Server]{
	ResourceNameSingular:   "Server",
	ShortDescriptionAdd:    "Add a label to a Server",
	ShortDescriptionRemove: "Remove a label from a Server",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Server().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Server().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.Server, error) {
		server, _, err := s.Client().Server().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if server == nil {
			return nil, fmt.Errorf("Server not found: %s", idOrName)
		}
		return server, nil
	},
	SetLabels: func(s state.State, server *hcloud.Server, labels map[string]string) error {
		opts := hcloud.ServerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Server().Update(s, server, opts)
		return err
	},
	GetLabels: func(server *hcloud.Server) map[string]string {
		return server.Labels
	},
	GetIDOrName: func(server *hcloud.Server) string {
		return strconv.FormatInt(server.ID, 10)
	},
	FetchBatch: func(s state.State, idOrNames []string) ([]*hcloud.Server, []error) {
		servers := make([]*hcloud.Server, len(idOrNames))
		errors := make([]error, len(idOrNames))

		var wg sync.WaitGroup
		for i, idOrName := range idOrNames {
			wg.Add(1)
			go func(idx int, id string) {
				defer wg.Done()
				server, _, err := s.Client().Server().Get(s, id)
				if err != nil {
					errors[idx] = err
					return
				}
				if server == nil {
					errors[idx] = fmt.Errorf("Server not found: %s", id)
					return
				}
				servers[idx] = server
			}(i, idOrName)
		}
		wg.Wait()

		return servers, errors
	},
	SetLabelsBatch: func(s state.State, servers []*hcloud.Server, labels map[string]string) []error {
		errors := make([]error, len(servers))

		var wg sync.WaitGroup
		for i, server := range servers {
			if server == nil {
				continue
			}

			wg.Add(1)
			go func(idx int, srv *hcloud.Server) {
				defer wg.Done()
				opts := hcloud.ServerUpdateOpts{
					Labels: labels,
				}
				_, _, err := s.Client().Server().Update(s, srv, opts)
				errors[idx] = err
			}(i, server)
		}
		wg.Wait()

		return errors
	},
}
