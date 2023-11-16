package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type PlacementGroupClient interface {
	hcloud.IPlacementGroupClient
	Names() []string
	LabelKeys(string) []string
}

func NewPlacementGroupClient(client hcloud.IPlacementGroupClient) PlacementGroupClient {
	return &placementGroupClient{
		IPlacementGroupClient: client,
	}
}

type placementGroupClient struct {
	hcloud.IPlacementGroupClient
}

func (c *placementGroupClient) Names() []string {
	placementGroups, err := c.All(context.Background())
	if err != nil || len(placementGroups) == 0 {
		return nil
	}
	names := make([]string, len(placementGroups))
	for i, firewall := range placementGroups {
		name := firewall.Name
		if name == "" {
			name = strconv.FormatInt(firewall.ID, 10)
		}
		names[i] = name
	}
	return names
}

func (c *placementGroupClient) LabelKeys(idOrName string) []string {
	placementGroups, _, err := c.Get(context.Background(), idOrName)
	if err != nil || placementGroups == nil || len(placementGroups.Labels) == 0 {
		return nil
	}
	return labelKeys(placementGroups.Labels)
}
