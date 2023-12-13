package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type FirewallClient interface {
	hcloud.IFirewallClient
	Names() []string
	LabelKeys(string) []string
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
func (c *firewallClient) Names() []string {
	firewalls, err := c.All(context.Background())
	if err != nil || len(firewalls) == 0 {
		return nil
	}
	names := make([]string, len(firewalls))
	for i, firewall := range firewalls {
		name := firewall.Name
		if name == "" {
			name = strconv.FormatInt(firewall.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the firewall with the passed idOrName.
func (c *firewallClient) LabelKeys(idOrName string) []string {
	firewall, _, err := c.Get(context.Background(), idOrName)
	if err != nil || firewall == nil || len(firewall.Labels) == 0 {
		return nil
	}
	return labelKeys(firewall.Labels)
}
