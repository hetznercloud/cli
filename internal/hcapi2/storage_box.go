package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxClient interface {
	hcloud.IStorageBoxClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
	SnapshotLabelKeys(context.Context, string, string) ([]string, error)
}

func NewStorageBoxClient(client hcloud.IStorageBoxClient) StorageBoxClient {
	return &storageBoxClient{
		IStorageBoxClient: client,
	}
}

type storageBoxClient struct {
	hcloud.IStorageBoxClient
}

// Names obtains a list of available Storage Boxes. It returns nil if Storage Box
// names could not be fetched or none are available.
func (c *storageBoxClient) Names(ctx context.Context) ([]string, error) {
	storageBoxes, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(storageBoxes, func(storageBox *hcloud.StorageBox) int64 { return storageBox.ID }, func(storageBox *hcloud.StorageBox) string { return storageBox.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels assigned to
// the Storage Box with the passed name or id.
func (c *storageBoxClient) LabelKeys(ctx context.Context, nameOrID string) ([]string, error) {
	storageBox, _, err := c.Get(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	if storageBox == nil {
		return nil, nil
	}
	return labelKeys(storageBox.Labels), nil
}

// SnapshotLabelKeys returns a slice containing the keys of all labels assigned to
// the Storage Box Snapshot with the passed name or id.
func (c *storageBoxClient) SnapshotLabelKeys(ctx context.Context, storageBoxNameOrID, snapshotNameOrID string) ([]string, error) {
	storageBox, _, err := c.Get(ctx, storageBoxNameOrID)
	if err != nil {
		return nil, err
	}
	if storageBox == nil {
		return nil, nil
	}
	storageBoxSnapshot, _, err := c.GetSnapshot(ctx, storageBox, snapshotNameOrID)
	if err != nil {
		return nil, err
	}
	if storageBoxSnapshot == nil {
		return nil, nil
	}
	return labelKeys(storageBoxSnapshot.Labels), nil
}
