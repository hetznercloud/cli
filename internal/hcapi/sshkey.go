package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// SSHKeyClient embeds the Hetzner Cloud SSHKey client and provides some
// additional helper functions.
type SSHKeyClient struct {
	*hcloud.SSHKeyClient
}

// SSHKeyNames obtains a list of available SSH keys. It returns nil if SSH key
// names could not be fetched or none are available.
func (c *SSHKeyClient) SSHKeyNames() []string {
	sshKeys, err := c.All(context.Background())
	if err != nil || len(sshKeys) == 0 {
		return nil
	}
	names := make([]string, len(sshKeys))
	for i, key := range sshKeys {
		name := key.Name
		if name == "" {
			name = strconv.Itoa(key.ID)
		}
		names[i] = name
	}
	return names
}

// SSHKeyLabelKeys returns a slice containing the keys of all labels
// assigned to the SSH Key with the passed idOrName.
func (c *SSHKeyClient) SSHKeyLabelKeys(idOrName string) []string {
	sshKey, _, err := c.Get(context.Background(), idOrName)
	if err != nil || sshKey == nil || len(sshKey.Labels) == 0 {
		return nil
	}
	return lkeys(sshKey.Labels)
}
