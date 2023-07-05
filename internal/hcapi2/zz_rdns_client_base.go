// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"net"
)

// RDNSClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/v2/hcloud.RDNSClient".
type RDNSClientBase interface {
	ChangeDNSPtr(context.Context, hcloud.RDNSSupporter, net.IP, *string) (*hcloud.Action, *hcloud.Response, error)
}
