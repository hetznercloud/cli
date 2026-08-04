package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type FirewallClient interface {
	hcloud.IFirewallClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

func NewFirewallClient(client hcloud.IFirewallClient) FirewallClient {
	return &firewallClient{
		IFirewallClient: client,
	}
}

// FirewallClient embeds the Hetzner Cloud Firewall client and provides
// some additional helper functions.
type firewallClient struct {
	hcloud.IFirewallClient
}

// Names obtains a list of available firewalls. It returns nil if
// the firewall names could not be fetched or there were no firewalls.
func (c *firewallClient) Names(ctx context.Context) ([]string, error) {
	firewalls, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(firewalls, func(firewall *hcloud.Firewall) int64 { return firewall.ID }, func(firewall *hcloud.Firewall) string { return firewall.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the firewall with the passed idOrName.
func (c *firewallClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	firewall, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if firewall == nil {
		return nil, nil
	}
	return labelKeys(firewall.Labels), nil
}
