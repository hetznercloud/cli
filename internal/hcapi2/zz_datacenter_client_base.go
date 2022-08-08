// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// DatacenterClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/hcloud.DatacenterClient".
type DatacenterClientBase interface {
	All(context.Context) ([]*hcloud.Datacenter, error)
	Get(context.Context, string) (*hcloud.Datacenter, *hcloud.Response, error)
	GetByID(context.Context, int) (*hcloud.Datacenter, *hcloud.Response, error)
	GetByName(context.Context, string) (*hcloud.Datacenter, *hcloud.Response, error)
	List(context.Context, hcloud.DatacenterListOpts) ([]*hcloud.Datacenter, *hcloud.Response, error)
}
