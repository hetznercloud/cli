package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// ISOClient embeds the Hetzner Cloud ISO client and provides some additional
// helper functions.
type ISOClient struct {
	*hcloud.ISOClient
}

// ISONames obtains a list of available ISOs for the current account. It
// returns nil if the current project has no ISOs or the ISO names could not be
// fetched.
func (c *ISOClient) ISONames() []string {
	isos, err := c.All(context.Background())
	if err != nil || len(isos) == 0 {
		return nil
	}
	names := make([]string, len(isos))
	for i, iso := range isos {
		name := iso.Name
		if name == "" {
			name = strconv.Itoa(iso.ID)
		}
		names[i] = name
	}
	return names
}
