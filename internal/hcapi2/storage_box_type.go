package hcapi2

import (
	"context"
	"sync"

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

type storageBoxTypeClient struct {
	hcloud.IStorageBoxTypeClient

	sbTypeByID map[int64]*hcloud.StorageBoxType
	once       sync.Once
	err        error
}

// Names returns a slice of all available storage box types.
func (c *storageBoxTypeClient) Names() []string {
	sts, err := c.All(context.Background())
	if err != nil || len(sts) == 0 {
		return nil
	}
	names := make([]string, len(sts))
	for i, st := range sts {
		names[i] = st.Name
	}
	return names
}

func (c *storageBoxTypeClient) init() error {
	c.once.Do(func() {
		storageBoxTypes, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(storageBoxTypes) == 0 {
			return
		}
		c.sbTypeByID = make(map[int64]*hcloud.StorageBoxType, len(storageBoxTypes))
		for _, sbt := range storageBoxTypes {
			c.sbTypeByID[sbt.ID] = sbt
		}
	})
	return c.err
}
