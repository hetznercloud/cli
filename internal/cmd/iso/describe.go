package iso

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a iso.
var DescribeCmd = base.DescribeCmd[*hcloud.ISO]{
	ResourceNameSingular: "iso",
	ShortDescription:     "Describe a iso",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Location().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.ISO, any, error) {
		iso, _, err := s.Client().ISO().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return iso, hcloud.SchemaFromISO(iso), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, iso *hcloud.ISO) error {
		cmd.Printf("ID:\t\t%d\n", iso.ID)
		cmd.Printf("Name:\t\t%s\n", iso.Name)
		cmd.Printf("Description:\t%s\n", iso.Description)
		cmd.Printf("Type:\t\t%s\n", iso.Type)
		cmd.Print(util.DescribeDeprecation(iso))

		architecture := "-"
		if iso.Architecture != nil {
			architecture = string(*iso.Architecture)
		}
		cmd.Printf("Architecture:\t%s\n", architecture)

		return nil
	},
}
