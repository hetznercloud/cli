package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ISOClient embeds the Hetzner Cloud iso client and provides some
// additional helper functions.
type ISOClient interface {
	hcloud.IISOClient
	Names() []string
}

func NewISOClient(client hcloud.IISOClient) ISOClient {
	return &isoClient{
		IISOClient: client,
	}
}

type isoClient struct {
	hcloud.IISOClient
}

// Names obtains a list of available data centers. It returns nil if
// iso names could not be fetched.
func (c *isoClient) Names() []string {
	isos, err := c.All(context.Background())
	if err != nil || len(isos) == 0 {
		return nil
	}
	names := make([]string, len(isos))
	for i, iso := range isos {
		name := iso.Name
		if name == "" {
			name = strconv.FormatInt(iso.ID, 10)
		}
		names[i] = name
	}
	return names
}
