package iso

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// DescribeCmd defines a command for describing a iso.
var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "iso",
	ShortDescription:     "Describe a iso",
	JSONKeyGetByID:       "iso",
	JSONKeyGetByName:     "isos",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.Location().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.ISO().Get(ctx, idOrName)
	},
	PrintText: func(ctx context.Context, client hcapi2.Client, resource interface{}) error {
		iso := resource.(*hcloud.ISO)

		fmt.Printf("ID:\t\t%d\n", iso.ID)
		fmt.Printf("Name:\t\t%s\n", iso.Name)
		fmt.Printf("Description:\t%s\n", iso.Description)
		fmt.Printf("Type:\t\t%s\n", iso.Type)
		return nil
	},
}
