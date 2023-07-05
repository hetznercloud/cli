// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// LocationClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/v2/hcloud.LocationClient".
type LocationClientBase interface {
	All(context.Context) ([]*hcloud.Location, error)
	Get(context.Context, string) (*hcloud.Location, *hcloud.Response, error)
	GetByID(context.Context, int) (*hcloud.Location, *hcloud.Response, error)
	GetByName(context.Context, string) (*hcloud.Location, *hcloud.Response, error)
	List(context.Context, hcloud.LocationListOpts) ([]*hcloud.Location, *hcloud.Response, error)
}
