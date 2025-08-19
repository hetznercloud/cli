package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxTypeClient interface {
	hcloud.IStorageBoxTypeClient
	Names() []string
}

func NewStorageBoxTypeClient(client hcloud.IStorageBoxTypeClient) StorageBoxTypeClient {
	return &storageBoxTypeClient{
		IStorageBoxTypeClient: client,
	}
}

// Names returns a slice of all available Storage Box types.
func (c *storageBoxTypeClient) Names() []string {
	types, err := c.All(context.Background())
	if err != nil || len(types) == 0 {
		return nil
	}
	names := make([]string, len(types))
	for i, t := range types {
		names[i] = t.Name
	}
	return names
}

type storageBoxTypeClient struct {
	hcloud.IStorageBoxTypeClient
}
