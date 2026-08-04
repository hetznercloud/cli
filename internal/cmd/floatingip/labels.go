package floatingip

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.FloatingIP]{
	ResourceNameSingular:   "Floating IP",
	ShortDescriptionAdd:    "Add a label to a Floating IP",
	ShortDescriptionRemove: "Remove a label from a Floating IP",
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.FloatingIP().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.FloatingIP().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.FloatingIP, error) {
		floatingIP, _, err := s.Client().FloatingIP().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if floatingIP == nil {
			return nil, fmt.Errorf("Floating IP not found: %s", idOrName)
		}
		return floatingIP, nil
	},
	SetLabels: func(ctx context.Context, s state.State, floatingIP *hcloud.FloatingIP, labels map[string]string) error {
		opts := hcloud.FloatingIPUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().FloatingIP().Update(ctx, floatingIP, opts)
		return err
	},
	GetLabels: func(floatingIP *hcloud.FloatingIP) map[string]string {
		return floatingIP.Labels
	},
	GetIDOrName: func(floatingIP *hcloud.FloatingIP) string {
		return strconv.FormatInt(floatingIP.ID, 10)
	},
}
