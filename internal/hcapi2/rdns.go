package hcapi2

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// NewRDNSClient embeds the Hetzner Cloud rdns client.
type RDNSClient interface {
	hcloud.IRDNSClient
}

func NewRDNSClient(client hcloud.IRDNSClient) RDNSClient {
	return &rdnsClient{
		IRDNSClient: client,
	}
}

type rdnsClient struct {
	hcloud.IRDNSClient
}
