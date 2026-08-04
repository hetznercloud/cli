package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxTypeClient interface {
	hcloud.IStorageBoxTypeClient
	Names(context.Context) ([]string, error)
}

func NewStorageBoxTypeClient(client hcloud.IStorageBoxTypeClient) StorageBoxTypeClient {
	return &storageBoxTypeClient{
		IStorageBoxTypeClient: client,
	}
}

type storageBoxTypeClient struct {
	hcloud.IStorageBoxTypeClient
}

// Names returns a slice of all available storage box types.
func (c *storageBoxTypeClient) Names(ctx context.Context) ([]string, error) {
	storageBoxTypes, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(storageBoxTypes))
	for i, st := range storageBoxTypes {
		names[i] = st.Name
	}
	return names, nil
}
