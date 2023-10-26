package server

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "server",
	ShortDescriptionAdd:    "Add a label to a server",
	ShortDescriptionRemove: "Remove a label from a server",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Server().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Server().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int64, error) {
		server, _, err := client.Server().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if server == nil {
			return nil, 0, fmt.Errorf("server not found: %s", idOrName)
		}
		return server.Labels, server.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int64, labels map[string]string) error {
		opts := hcloud.ServerUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.Server().Update(ctx, &hcloud.Server{ID: id}, opts)
		return err
	},
}
