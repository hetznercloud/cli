package hcapi2

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// PricingClient embeds the Hetzner Cloud Pricing client and provides some
// additional helper functions.
type PricingClient interface {
	hcloud.IPricingClient
}

func NewPricingClient(client hcloud.IPricingClient) PricingClient {
	return &pricingClient{
		IPricingClient: client,
	}
}

type pricingClient struct {
	hcloud.IPricingClient
}
