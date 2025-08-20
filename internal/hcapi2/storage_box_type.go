package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxTypeClient interface {
	hcloud.IStorageBoxTypeClient
	Names() []string
	StorageBoxTypeName(id int64) string
	StorageBoxTypeDescription(id int64) string
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

// StorageBoxTypeName obtains the name of the storage box type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *storageBoxTypeClient) StorageBoxTypeName(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	storageBoxType, ok := c.sbTypeByID[id]
	if !ok || storageBoxType.Name == "" {
		return strconv.FormatInt(id, 10)
	}
	return storageBoxType.Name
}

// StorageBoxTypeDescription obtains the description of the storage box type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *storageBoxTypeClient) StorageBoxTypeDescription(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	storageBoxType, ok := c.sbTypeByID[id]
	if !ok || storageBoxType.Description == "" {
		return strconv.FormatInt(id, 10)
	}
	return storageBoxType.Description
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
