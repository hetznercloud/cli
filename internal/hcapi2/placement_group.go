package hcapi2

import (
	"context"
	"strconv"
)

type PlacementGroupClient interface {
	PlacementGroupClientBase
	Names() []string
	LabelKeys(string) []string
}

func NewPlacementGroupClient(client PlacementGroupClientBase) PlacementGroupClient {
	return &placementGroupClient{
		PlacementGroupClientBase: client,
	}
}

type placementGroupClient struct {
	PlacementGroupClientBase
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
			name = strconv.Itoa(firewall.ID)
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
