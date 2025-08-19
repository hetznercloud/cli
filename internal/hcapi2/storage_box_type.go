package hcapi2

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxTypeClient interface {
	hcloud.IStorageBoxTypeClient
}

func NewStorageBoxTypeClient(client hcloud.IStorageBoxTypeClient) StorageBoxTypeClient {
	return &storageBoxTypeClient{
		IStorageBoxTypeClient: client,
	}
}

type storageBoxTypeClient struct {
	hcloud.IStorageBoxTypeClient
}
