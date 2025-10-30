package iso

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// DescribeCmd defines a command for describing a iso.
var DescribeCmd = base.DescribeCmd[*hcloud.ISO]{
	ResourceNameSingular: "ISO",
	ShortDescription:     "Describe an ISO",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.ISO().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.ISO, any, error) {
		iso, _, err := s.Client().ISO().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return iso, hcloud.SchemaFromISO(iso), nil
	},
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, iso *hcloud.ISO) error {
		fmt.Fprint(out, DescribeISO(iso))
		return nil
	},
}

func DescribeISO(iso *hcloud.ISO) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ID:\t%d\n", iso.ID)
	fmt.Fprintf(&sb, "Name:\t%s\n", iso.Name)
	fmt.Fprintf(&sb, "Description:\t%s\n", iso.Description)
	fmt.Fprintf(&sb, "Type:\t%s\n", iso.Type)
	fmt.Fprintf(&sb, "%s", util.DescribeDeprecation(iso))

	architecture := "-"
	if iso.Architecture != nil {
		architecture = string(*iso.Architecture)
	}
	fmt.Fprintf(&sb, "Architecture:\t%s\n", architecture)

	return sb.String()
}
