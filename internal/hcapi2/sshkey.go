package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// SSHKeyClient embeds the Hetzner Cloud SSHKey client and provides some
// additional helper functions.
type SSHKeyClient interface {
	hcloud.ISSHKeyClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
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
func (c *sshKeyClient) Names(ctx context.Context) ([]string, error) {
	sshKeys, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(sshKeys, func(key *hcloud.SSHKey) int64 { return key.ID }, func(key *hcloud.SSHKey) string { return key.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the SSH Key with the passed idOrName.
func (c *sshKeyClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	sshKey, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if sshKey == nil {
		return nil, nil
	}
	return labelKeys(sshKey.Labels), nil
}
