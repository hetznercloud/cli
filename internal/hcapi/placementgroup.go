package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type PlacementGroupClient struct {
	*hcloud.PlacementGroupClient
}

func (c *PlacementGroupClient) PlacementGroupNames() []string {
	placementGroups, err := c.All(context.Background())
	if err != nil || len(placementGroups) == 0 {
		return nil
	}
	names := make([]string, len(placementGroups))
	for i, firewall := range placementGroups {
		name := firewall.Name
		if name == "" {
			name = strconv.Itoa(firewall.ID)
		}
		names[i] = name
	}
	return names
}

func (c *PlacementGroupClient) PlacementGroupLabelKeys(idOrName string) []string {
	placementGroups, _, err := c.Get(context.Background(), idOrName)
	if err != nil || placementGroups == nil || len(placementGroups.Labels) == 0 {
		return nil
	}
	return lkeys(placementGroups.Labels)
}
