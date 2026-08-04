package primaryip

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.PrimaryIP]{
	ResourceNameSingular:   "Primary IP",
	ShortDescriptionAdd:    "Add a label to a Primary IP",
	ShortDescriptionRemove: "Remove a label from a Primary IP",
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.PrimaryIP().Names(false, false, nil) },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.PrimaryIP().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.PrimaryIP, error) {
		primaryIP, _, err := s.Client().PrimaryIP().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if primaryIP == nil {
			return nil, fmt.Errorf("primaryIP not found: %s", idOrName)
		}
		return primaryIP, nil
	},
	SetLabels: func(ctx context.Context, s state.State, primaryIP *hcloud.PrimaryIP, labels map[string]string) error {
		opts := hcloud.PrimaryIPUpdateOpts{
			Labels: &labels,
		}
		_, _, err := s.Client().PrimaryIP().Update(ctx, primaryIP, opts)
		return err
	},
	GetLabels: func(primaryIP *hcloud.PrimaryIP) map[string]string {
		return primaryIP.Labels
	},
	GetIDOrName: func(primaryIP *hcloud.PrimaryIP) string {
		return strconv.FormatInt(primaryIP.ID, 10)
	},
}
