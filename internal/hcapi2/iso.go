package hcapi2

import (
	"context"
	"strconv"
)

// ISOClient embeds the Hetzner Cloud iso client and provides some
// additional helper functions.
type ISOClient interface {
	ISOClientBase
	Names() []string
}

func NewISOClient(client ISOClientBase) ISOClient {
	return &isoClient{
		ISOClientBase: client,
	}
}

type isoClient struct {
	ISOClientBase
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
