package hcapi2

import (
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxClient interface {
	hcloud.IStorageBoxClient
}

func NewStorageBoxClient(client hcloud.IStorageBoxClient) StorageBoxClient {
	return &storageBoxClient{
		IStorageBoxClient: client,
	}
}

type storageBoxClient struct {
	hcloud.IStorageBoxClient

	once sync.Once
	err  error
}
