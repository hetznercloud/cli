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
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		id, err := strconv.ParseInt(idOrName, 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid snapshot or backup ID %q", idOrName)
		}
		image, _, err := s.Client().Image().GetByID(s, id)
		if err != nil {
			return nil, 0, err
		}
		if image == nil {
			return nil, 0, fmt.Errorf("image not found: %s", idOrName)
		}
		return image.Labels, image.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.ImageUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Image().Update(s, &hcloud.Image{ID: id}, opts)
		return err
	},
}
