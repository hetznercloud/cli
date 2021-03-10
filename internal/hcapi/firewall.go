package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// FirewallClient embeds the Hetzner Cloud Firewall client and provides
// some additional helper functions.
type FirewallClient struct {
	*hcloud.FirewallClient
}

// FirewallNames obtains a list of available firewalls. It returns nil if
// the firewall names could not be fetched or there were no firewalls.
func (c *FirewallClient) FirewallNames() []string {
	firewalls, err := c.All(context.Background())
	if err != nil || len(firewalls) == 0 {
		return nil
	}
	names := make([]string, len(firewalls))
	for i, firewall := range firewalls {
		name := firewall.Name
		if name == "" {
			name = strconv.Itoa(firewall.ID)
		}
		names[i] = name
	}
	return names
}

// FirewallLabelKeys returns a slice containing the keys of all labels
// assigned to the firewall with the passed idOrName.
func (c *FirewallClient) FirewallLabelKeys(idOrName string) []string {
	firewall, _, err := c.Get(context.Background(), idOrName)
	if err != nil || firewall == nil || len(firewall.Labels) == 0 {
		return nil
	}
	return lkeys(firewall.Labels)
}
