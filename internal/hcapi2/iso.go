package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ISOClient embeds the Hetzner Cloud iso client and provides some
// additional helper functions.
type ISOClient interface {
	hcloud.IISOClient
	Names(context.Context) ([]string, error)
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
func (c *isoClient) Names(ctx context.Context) ([]string, error) {
	isos, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(isos, func(iso *hcloud.ISO) int64 { return iso.ID }, func(iso *hcloud.ISO) string { return iso.Name }), nil
}
