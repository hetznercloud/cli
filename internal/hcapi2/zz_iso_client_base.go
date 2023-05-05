// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// ISOClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/hcloud.ISOClient".
type ISOClientBase interface {
	All(context.Context) ([]*hcloud.ISO, error)
	AllWithOpts(context.Context, hcloud.ISOListOpts) ([]*hcloud.ISO, error)
	Get(context.Context, string) (*hcloud.ISO, *hcloud.Response, error)
	GetByID(context.Context, int) (*hcloud.ISO, *hcloud.Response, error)
	GetByName(context.Context, string) (*hcloud.ISO, *hcloud.Response, error)
	List(context.Context, hcloud.ISOListOpts) ([]*hcloud.ISO, *hcloud.Response, error)
}
