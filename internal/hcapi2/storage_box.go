package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type StorageBoxClient interface {
	hcloud.IStorageBoxClient
	Names() []string
	LabelKeys(string) []string
	SnapshotLabelKeys(string, string) []string
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
func (c *storageBoxClient) Names() []string {
	storageBoxes, err := c.All(context.Background())
	if err != nil || len(storageBoxes) == 0 {
		return nil
	}
	names := make([]string, len(storageBoxes))
	for i, storageBox := range storageBoxes {
		name := storageBox.Name
		if name == "" {
			name = strconv.FormatInt(storageBox.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels assigned to
// the Storage Box with the passed name or id.
func (c *storageBoxClient) LabelKeys(nameOrID string) []string {
	storageBox, _, err := c.Get(context.Background(), nameOrID)
	if err != nil || storageBox == nil || len(storageBox.Labels) == 0 {
		return nil
	}
	return labelKeys(storageBox.Labels)
}

// SnapshotLabelKeys returns a slice containing the keys of all labels assigned to
// the Storage Box Snapshot with the passed name or id.
func (c *storageBoxClient) SnapshotLabelKeys(storageBoxNameOrID, snapshotNameOrID string) []string {
	storageBox, _, err := c.Get(context.Background(), storageBoxNameOrID)
	if err != nil || storageBox == nil {
		return nil
	}
	storageBoxSnapshot, _, err := c.GetSnapshot(context.Background(), storageBox, snapshotNameOrID)
	if err != nil || storageBoxSnapshot == nil || len(storageBoxSnapshot.Labels) == 0 {
		return nil
	}
	return labelKeys(storageBoxSnapshot.Labels)
}
