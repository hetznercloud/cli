package image

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "image",
	ShortDescriptionAdd:    "Add a label to an image",
	ShortDescriptionRemove: "Remove a label from an image",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Image().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Image().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		id, err := strconv.ParseInt(idOrName, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid snapshot or backup ID %q", idOrName)
		}
		image, _, err := s.Client().Image().GetByID(s, id)
		if err != nil {
			return nil, err
		}
		if image == nil {
			return nil, fmt.Errorf("image not found: %s", idOrName)
		}
		return image, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		image := resource.(*hcloud.Image)
		opts := hcloud.ImageUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Image().Update(s, image, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		image := resource.(*hcloud.Image)
		return image.Labels
	},
	GetIDOrName: func(resource any) string {
		image := resource.(*hcloud.Image)
		return strconv.FormatInt(image.ID, 10)
	},
}
