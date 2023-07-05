// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// PlacementGroupClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/v2/hcloud.PlacementGroupClient".
type PlacementGroupClientBase interface {
	All(context.Context) ([]*hcloud.PlacementGroup, error)
	AllWithOpts(context.Context, hcloud.PlacementGroupListOpts) ([]*hcloud.PlacementGroup, error)
	Create(context.Context, hcloud.PlacementGroupCreateOpts) (hcloud.PlacementGroupCreateResult, *hcloud.Response, error)
	Delete(context.Context, *hcloud.PlacementGroup) (*hcloud.Response, error)
	Get(context.Context, string) (*hcloud.PlacementGroup, *hcloud.Response, error)
	GetByID(context.Context, int) (*hcloud.PlacementGroup, *hcloud.Response, error)
	GetByName(context.Context, string) (*hcloud.PlacementGroup, *hcloud.Response, error)
	List(context.Context, hcloud.PlacementGroupListOpts) ([]*hcloud.PlacementGroup, *hcloud.Response, error)
	Update(context.Context, *hcloud.PlacementGroup, hcloud.PlacementGroupUpdateOpts) (*hcloud.PlacementGroup, *hcloud.Response, error)
}
