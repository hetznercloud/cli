package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// SSHKeyClient embeds the Hetzner Cloud SSHKey client and provides some
// additional helper functions.
type SSHKeyClient interface {
	hcloud.ISSHKeyClient
	Names() []string
	LabelKeys(idOrName string) []string
}

func NewSSHKeyClient(client hcloud.ISSHKeyClient) SSHKeyClient {
	return &sshKeyClient{
		ISSHKeyClient: client,
	}
}

type sshKeyClient struct {
	hcloud.ISSHKeyClient
}

// Names obtains a list of available SSH keys. It returns nil if SSH key
// names could not be fetched or none are available.
func (c *sshKeyClient) Names() []string {
	sshKeys, err := c.All(context.Background())
	if err != nil || len(sshKeys) == 0 {
		return nil
	}
	names := make([]string, len(sshKeys))
	for i, key := range sshKeys {
		name := key.Name
		if name == "" {
			name = strconv.FormatInt(key.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the SSH Key with the passed idOrName.
func (c *sshKeyClient) LabelKeys(idOrName string) []string {
	sshKey, _, err := c.Get(context.Background(), idOrName)
	if err != nil || sshKey == nil || len(sshKey.Labels) == 0 {
		return nil
	}
	return labelKeys(sshKey.Labels)
}
