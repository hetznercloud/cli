package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type PlacementGroupClient interface {
	hcloud.IPlacementGroupClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

func NewPlacementGroupClient(client hcloud.IPlacementGroupClient) PlacementGroupClient {
	return &placementGroupClient{
		IPlacementGroupClient: client,
	}
}

type placementGroupClient struct {
	hcloud.IPlacementGroupClient
}

func (c *placementGroupClient) Names(ctx context.Context) ([]string, error) {
	placementGroups, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(placementGroups, func(group *hcloud.PlacementGroup) int64 { return group.ID }, func(group *hcloud.PlacementGroup) string { return group.Name }), nil
}

func (c *placementGroupClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	placementGroup, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if placementGroup == nil {
		return nil, nil
	}
	return labelKeys(placementGroup.Labels), nil
}
